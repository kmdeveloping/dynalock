package dynalock

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/kmdeveloping/dynalock/internal/lock"
	"github.com/kmdeveloping/dynalock/internal/models"
	"github.com/kmdeveloping/dynalock/providers"
)

const (
	defaultPartitionKeyName = "key"
)

var (
	defaultOwnerName = uuid.NewString()
)

type Client struct {
	TableName        string
	OwnerName        string
	PartitionKeyName string
	EncryptionKey    string
	Encryption       bool
	Lock             *lock.LockManager
	DynamoDB         providers.DynamoDbProvider
}

type ClientOption func(*Client)

func NewDynalockClient(dynamoDb providers.DynamoDbProvider, tableName string, opts ...ClientOption) *Client {

	client := &Client{
		TableName:        tableName,
		OwnerName:        defaultOwnerName,
		PartitionKeyName: defaultPartitionKeyName,
		DynamoDB:         dynamoDb,
	}

	for _, opt := range opts {
		opt(client)
	}

	lockManagerOptions := lock.LockManagerOptions{
		PartitionKeyName: client.PartitionKeyName,
		TableName:        client.TableName,
		OwnerName:        client.OwnerName,
		WithEncryption:   client.Encryption,
	}

	client.Lock = lock.NewLockManager(dynamoDb, lockManagerOptions)

	return client
}

func (c *Client) CreateTable(ctx context.Context, tableName string, opts ...models.CreateTableOption) (*dynamodb.CreateTableOutput, error) {
	createTableOptions := &models.CreateDynamoDBTableOptions{
		TableName:    tableName,
		BillingMode:  "PAY_PER_REQUEST",
		PartitionKey: c.PartitionKeyName,
	}
	for _, opt := range opts {
		opt(createTableOptions)
	}

	return c.Lock.CreateLockTable(ctx, createTableOptions)
}

func (c *Client) AquireLock(ctx context.Context, key string, opts ...models.AcquireLockOption) (*models.Lock, error) {
	if err := c.Lock.CanAcquireLock(ctx, key); err != nil {
		return nil, err
	}

	acquireLockOptions := &models.AcquireLockOptions{
		PartitionKey: key,
	}

	for _, opt := range opts {
		opt(acquireLockOptions)
	}

	return c.Lock.AcquireLock(ctx, key, acquireLockOptions)
}
