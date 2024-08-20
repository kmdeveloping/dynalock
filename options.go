package dynalock

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kmdeveloping/dynalock/internal/models"
)

func WithPartitionKeyName(pkn string) ClientOption {
	return func(client *Client) {
		client.PartitionKeyName = pkn
	}
}

func WithLockOwnerName(name string) ClientOption {
	return func(client *Client) {
		client.OwnerName = name
	}
}

func WithDataEncryption(key string) EncryptedClientOption {
	return func(client *Client) {
		client.EncryptionKey = key
		client.Encryption = true
	}
}

func WithData(data []byte) models.AcquireLockOption {
	return func(opt *models.AcquireLockOptions) {
		opt.Data = data
	}
}

func DeleteOnRelease() models.AcquireLockOption {
	return func(opt *models.AcquireLockOptions) {
		opt.DeleteOnRelease = true
	}
}

func ReplaceData() models.AcquireLockOption {
	return func(opt *models.AcquireLockOptions) {
		opt.ReplaceData = true
	}
}

func FailIfLocked() models.AcquireLockOption {
	return func(opt *models.AcquireLockOptions) {
		opt.FailIfLocked = true
	}
}
func WithAdditionalAttributes(attr map[string]types.AttributeValue) models.AcquireLockOption {
	return func(opt *models.AcquireLockOptions) {
		opt.AdditionalAttributes = attr
	}
}

func WithCustomPartitionKey(partitionKey string) models.CreateTableOption {
	return func(opt *models.CreateDynamoDBTableOptions) {
		opt.PartitionKey = partitionKey
	}
}

func WithProvisionedThroughput(provisionedThroughput *types.ProvisionedThroughput) models.CreateTableOption {
	return func(opt *models.CreateDynamoDBTableOptions) {
		opt.BillingMode = types.BillingModeProvisioned
		opt.ProvisionedThroughput = provisionedThroughput
	}
}
