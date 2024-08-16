package lock

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kmdeveloping/dynalock/internal/models"
	"github.com/kmdeveloping/dynalock/providers"
)

type LockManager struct {
	WithEncryption bool
	tableName      string
	dynamoClient   providers.DynamoDbProvider
}

func NewLockManager(dynamoClient providers.DynamoDbProvider, tableName string, withEncryption bool) *LockManager {
	return &LockManager{
		WithEncryption: withEncryption,
		tableName:      tableName,
		dynamoClient:   dynamoClient,
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
	return nil
}
