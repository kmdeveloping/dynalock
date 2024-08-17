package lock

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kmdeveloping/dynalock/internal/models"
	"github.com/kmdeveloping/dynalock/providers"
	"log"
)

type LockManager struct {
	WithEncryption   bool
	partitionKeyName string
	tableName        string
	ownerName        string
	dynamoClient     providers.DynamoDbProvider
}

type LockManagerOptions struct {
	WithEncryption   bool
	PartitionKeyName string
	TableName        string
	OwnerName        string
}

func NewLockManager(dc providers.DynamoDbProvider, opts LockManagerOptions) *LockManager {
	return &LockManager{
		WithEncryption:   opts.WithEncryption,
		partitionKeyName: opts.PartitionKeyName,
		tableName:        opts.TableName,
		ownerName:        opts.OwnerName,
		dynamoClient:     dc,
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
func (lm *LockManager) AcquireLock(ctx context.Context, key string, opt *models.AcquireLockOptions) (*models.Lock, error) {
	return nil, nil
}
func (lm *LockManager) ReleaseLock(ctx context.Context, key string) (bool, error) {
	return false, nil
}
