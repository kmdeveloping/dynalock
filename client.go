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
		OwnerName:        uuid.NewString(),
		DynamoDB:         dynamoDb,
		PartitionKeyName: defaultPartitionKeyName,
	}

	for _, opt := range opts {
		opt(client)
	}

	client.Lock = lock.NewLockManager(dynamoDb, tableName, client.Encryption)

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
