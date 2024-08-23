package lock

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kmdeveloping/dynalock/internal/models"
	. "github.com/kmdeveloping/dynalock/providers"
	"github.com/kmdeveloping/encrypticon"
	"log"
	"time"
)

type LockManager struct {
	withEncryption    bool
	partitionKeyName  string
	tableName         string
	ownerName         string
	encryptionManager *encrypticon.EncryptManager
	dynamoClient      DynamoDbProvider
}

type LockManagerOptions struct {
	WithEncryption   bool
	PartitionKeyName string
	TableName        string
	OwnerName        string
}

type LockManagerOptionsWithEncryption struct {
	LockManagerOptions
	EncryptionKey string
}

var _ LockManagerProvider = (*LockManager)(nil)

func NewLockManager(dc DynamoDbProvider, opts LockManagerOptions) *LockManager {
	return &LockManager{
		withEncryption:   opts.WithEncryption,
		partitionKeyName: opts.PartitionKeyName,
		tableName:        opts.TableName,
		ownerName:        opts.OwnerName,
		dynamoClient:     dc,
	}
}

func NewLockManagerWithEncryption(dc DynamoDbProvider, opts LockManagerOptionsWithEncryption) *LockManager {
	return &LockManager{
		withEncryption:    opts.WithEncryption,
		encryptionManager: encrypticon.NewEncryptManager(opts.EncryptionKey),
		partitionKeyName:  opts.PartitionKeyName,
		tableName:         opts.TableName,
		ownerName:         opts.OwnerName,
		dynamoClient:      dc,
	}
}

func (lm *LockManager) CreateLockTable(ctx context.Context, opt *models.CreateDynamoDBTableOptions) (*dynamodb.CreateTableOutput, error) {
	keySchema := []types.KeySchemaElement{
		{
			AttributeName: aws.String(opt.PartitionKey),
			KeyType:       types.KeyTypeHash,
		},
	}

	attributeDefinitions := []types.AttributeDefinition{
		{
			AttributeName: aws.String(opt.PartitionKey),
			AttributeType: types.ScalarAttributeTypeS,
		},
	}

	createTableInput := &dynamodb.CreateTableInput{
		TableName:            aws.String(lm.tableName),
		KeySchema:            keySchema,
		BillingMode:          opt.BillingMode,
		AttributeDefinitions: attributeDefinitions,
	}

	if opt.ProvisionedThroughput != nil {
		createTableInput.ProvisionedThroughput = opt.ProvisionedThroughput
	}

	return lm.dynamoClient.CreateTable(ctx, createTableInput)
}

func (lm *LockManager) CanAcquireLock(ctx context.Context, key string) error {
	query := &dynamodb.QueryInput{
		TableName:              aws.String(lm.tableName),
		KeyConditionExpression: aws.String("#key = :lock_key"),
		ExpressionAttributeNames: map[string]string{
			"#key": "key",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":lock_key": &types.AttributeValueMemberS{Value: key},
		},
	}

	q, err := lm.dynamoClient.Query(ctx, query)
	if err != nil {
		log.Printf("failed to query: %v", err)
		return err
	}

	if q.Items == nil {
		return nil
	}

	var locks []models.Lock
	locks, err = lm.unmarshalListLocks(q.Items)
	if err != nil {
		log.Printf("failed to unmarshal locks: %v", err)
		return err
	}

	for _, lock := range locks {
		if lock.Owner != lm.ownerName && !lock.IsReleased {
			return errors.New("lock is in use by another and not released")
		}
	}

	return nil
}
func (lm *LockManager) AcquireLock(ctx context.Context, opt *models.AcquireLockOptions) (*models.Lock, error) {
	lock := models.Lock{
		PartitionKey:    opt.PartitionKey,
		Owner:           lm.ownerName,
		Timestamp:       time.Now().Unix(),
		Ttl:             30 * 24 * 60 * 60,
		DeleteOnRelease: opt.DeleteOnRelease,
		IsReleased:      false,
		Data:            opt.Data,
	}

	if lm.withEncryption {
		encryptedString := lm.encryptionManager.Encrypt(string(opt.Data))
		lock.Data = []byte(encryptedString)
	}

	item, err := lm.marshalLockItem(lock)
	if err != nil {
		return nil, err
	}

	req := &dynamodb.PutItemInput{
		TableName: aws.String(lm.tableName),
		Item:      item,
	}

	_, err = lm.dynamoClient.PutItem(ctx, req)
	if err != nil {
		return nil, err
	}

	return &lock, nil
}
func (lm *LockManager) ReleaseLock(ctx context.Context, key string) (bool, error) {
	getInput := &dynamodb.GetItemInput{
		TableName: aws.String(lm.tableName),
		Key: map[string]types.AttributeValue{
			lm.partitionKeyName: &types.AttributeValueMemberS{Value: key},
		},
	}

	getOutput, err := lm.dynamoClient.GetItem(ctx, getInput)
	if err != nil {
		return false, err
	}

	if getOutput.Item == nil {
		return true, errors.New(fmt.Sprintf("no lock found for %s", key))
	}

	lockItem, err := lm.unmarshalLockItem(getOutput.Item)
	if err != nil {
		return false, err
	}

	if !lockItem.DeleteOnRelease {
		updateInput := &dynamodb.UpdateItemInput{
			TableName: aws.String(lm.tableName),
			Key: map[string]types.AttributeValue{
				lm.partitionKeyName: &types.AttributeValueMemberS{Value: key},
			},
			UpdateExpression: aws.String("SET #ts = :timestamp, isReleased = :isReleased"),
			ExpressionAttributeNames: map[string]string{
				"#ts":    attrTimestamp,
				"#owner": attrOwnerName,
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":owner":      &types.AttributeValueMemberS{Value: lm.ownerName},
				":timestamp":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", time.Now().Unix())},
				":isReleased": &types.AttributeValueMemberBOOL{Value: true},
			},
			ConditionExpression: aws.String("#owner = :owner"),
		}

		_, err = lm.dynamoClient.UpdateItem(ctx, updateInput)
		if err != nil {
			return false, fmt.Errorf("failed to update lock: %v", err)
		}

		return true, nil
	}

	return lm.deleteLockOnRelease(ctx, key)
}

func (lm *LockManager) deleteLockOnRelease(ctx context.Context, key string) (bool, error) {
	deleteInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(lm.tableName),
		Key: map[string]types.AttributeValue{
			lm.partitionKeyName: &types.AttributeValueMemberS{Value: key},
		},
		ExpressionAttributeNames: map[string]string{
			"#owner": attrOwnerName,
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":owner":  &types.AttributeValueMemberS{Value: lm.ownerName},
			":delete": &types.AttributeValueMemberBOOL{Value: true},
		},
		ConditionExpression: aws.String("#owner = :owner and deleteOnRelease = :delete"),
	}

	_, err := lm.dynamoClient.DeleteItem(ctx, deleteInput)
	if err != nil {
		return false, fmt.Errorf("failed to delete lock: %v", err)
	}

	return true, nil
}
