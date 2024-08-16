package models

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

type Lock struct {
	PartitionKey    string `dynamodbav:"key"`
	Owner           string `dynamodbav:"owner"`
	Timestamp       int64  `dynamodbav:"timestamp"`
	Ttl             int64  `dynamodbav:"ttl"`
	DeleteOnRelease bool   `dynamodbav:"deleteOnRelease"`
	IsReleased      bool   `dynamodbav:"isReleased"`
	Data            []byte `dynamodbav:"data"`
}

type CreateDynamoDBTableOptions struct {
	BillingMode           types.BillingMode
	ProvisionedThroughput *types.ProvisionedThroughput
	TableName             string
	PartitionKey          string
}

type AcquireLockOptions struct {
	PartitionKey         string
	Data                 []byte
	ReplaceData          bool
	DeleteOnRelease      bool
	FailIfLocked         bool
	AdditionalAttributes map[string]types.AttributeValue
}

type getLockOptions struct {
	partitionKeyName     string
	deleteLockOnRelease  bool
	replaceData          bool
	data                 []byte
	additionalAttributes map[string]types.AttributeValue
	failIfLocked         bool
}

type CreateTableOption func(*CreateDynamoDBTableOptions)
type AcquireLockOption func(*AcquireLockOptions)
