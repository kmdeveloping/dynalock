package dynalock_tests

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kmdeveloping/dynalock/internal/lock"
	"github.com/kmdeveloping/dynalock/internal/mocks"
	"github.com/kmdeveloping/dynalock/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type lock_manager_UnitTestSuite struct {
	suite.Suite

	lockManager           *lock.LockManager
	mockDynamoDb          *mocks.DynamoDbProvider
	mockContext           context.Context
	lockManagerOptions    lock.LockManagerOptions
	defaultTableName      string
	defaultLockOwner      string
	defaultPartitionKey   string
	deleteOnReleaseOption bool
	isReleasedOption      bool
	dataOption            []byte
	lockOutputValue       *models.Lock
}

func Test_lock_manager_UnitTestSuite(t *testing.T) {
	suite.Run(t, &lock_manager_UnitTestSuite{})
}

func (suite *lock_manager_UnitTestSuite) SetupSuite() {
	suite.defaultTableName = "mock_lock_table"
	suite.defaultPartitionKey = "mock_partition_key"
	suite.defaultLockOwner = "mock_lock_owner"
	suite.dataOption = []byte(`{"test": "mock_data"}`)
	suite.deleteOnReleaseOption = false
	suite.isReleasedOption = false
}

func (suite *lock_manager_UnitTestSuite) setupMocks() {
	suite.mockContext = context.Background()
	suite.mockDynamoDb = mocks.NewDynamoDbProvider(suite.T())

	suite.lockManagerOptions = lock.LockManagerOptions{
		WithEncryption:   false,
		PartitionKeyName: suite.defaultPartitionKey,
		TableName:        suite.defaultTableName,
		OwnerName:        suite.defaultLockOwner,
	}

	suite.lockOutputValue = &models.Lock{
		PartitionKey:    suite.defaultPartitionKey,
		Owner:           suite.defaultLockOwner,
		Timestamp:       time.Now().Unix(),
		Ttl:             5 * 60 * 60 * 24,
		DeleteOnRelease: suite.deleteOnReleaseOption,
		IsReleased:      suite.isReleasedOption,
		Data:            suite.dataOption,
	}

	suite.lockManager = lock.NewLockManager(suite.mockDynamoDb, suite.lockManagerOptions)
}

func (suite *lock_manager_UnitTestSuite) Test_CanAcquireLock_Success() {
	suite.setupMocks()
	suite.mockDynamoDb.EXPECT().
		Query(suite.mockContext, mock.Anything).
		Return(&dynamodb.QueryOutput{Count: 0, Items: nil}, nil).
		Once()

	err := suite.lockManager.CanAcquireLock(suite.mockContext, suite.defaultPartitionKey)
	suite.NoError(err)
}

func (suite *lock_manager_UnitTestSuite) Test_AcquireLock_Success() {
	suite.setupMocks()
	suite.mockDynamoDb.EXPECT().
		PutItem(suite.mockContext, mock.Anything).
		Return(&dynamodb.PutItemOutput{Attributes: suite.getAttributeValueMap(suite.defaultLockOwner)[0]}, nil).
		Once()

	options := &models.AcquireLockOptions{
		PartitionKey: suite.defaultPartitionKey,
		Data:         suite.dataOption,
	}

	lk, err := suite.lockManager.AcquireLock(suite.mockContext, options)
	suite.Nil(err)
	suite.NotNil(lk)
}

func (suite *lock_manager_UnitTestSuite) getAttributeValueMap(ownerName string) []map[string]types.AttributeValue {
	return []map[string]types.AttributeValue{
		{
			"key":             &types.AttributeValueMemberS{Value: suite.defaultPartitionKey},
			"owner":           &types.AttributeValueMemberS{Value: ownerName},
			"timestamp":       &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", suite.lockOutputValue.Timestamp)},
			"ttl":             &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", suite.lockOutputValue.Ttl)},
			"isReleased":      &types.AttributeValueMemberBOOL{Value: suite.isReleasedOption},
			"deleteOnRelease": &types.AttributeValueMemberBOOL{Value: suite.deleteOnReleaseOption},
		},
	}
}
