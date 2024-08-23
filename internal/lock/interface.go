package lock

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/kmdeveloping/dynalock/internal/models"
)

type LockManagerProvider interface {
	CreateLockTable(ctx context.Context, opt *models.CreateDynamoDBTableOptions) (*dynamodb.CreateTableOutput, error)
	CanAcquireLock(ctx context.Context, key string) error
	AcquireLock(ctx context.Context, opt *models.AcquireLockOptions) (*models.Lock, error)
	ReleaseLock(ctx context.Context, key string) (bool, error)
}
