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
	"log"
	"testing"
	"time"
)

type lock_manager_UnitTestSuite struct {
	suite.Suite

	lockManager             *lock.LockManager
	mockDynamoDb            *mocks.DynamoDbProvider
	mockContext             context.Context
	lockManagerOptions      lock.LockManagerOptions
	defaultPartitionKeyName string
	defaultTableName        string
	defaultLockOwner        string
	defaultPartitionKey     string
	deleteOnReleaseOption   bool
	isReleasedOption        bool
	dataOption              []byte
	lockOutputValue         *models.Lock

	useEncryptionLockManager bool
	encryptionKeyString      string
}

func Test_lock_manager_UnitTestSuite(t *testing.T) {
	suite.Run(t, &lock_manager_UnitTestSuite{})
}

func (suite *lock_manager_UnitTestSuite) SetupSuite() {
	suite.defaultPartitionKeyName = "key"
	suite.defaultTableName = "mock_lock_table"
	suite.defaultPartitionKey = "mock_partition_key"
	suite.defaultLockOwner = "mock_lock_owner"
	suite.dataOption = []byte(`{"test": "mock_data"}`)
	suite.deleteOnReleaseOption = false
	suite.isReleasedOption = false
	suite.useEncryptionLockManager = false
	suite.encryptionKeyString = "buh3c64sldpse2yxv5ujzyaa6pd7sdns"
}

func (suite *lock_manager_UnitTestSuite) setupMocks() {
	suite.mockContext = context.Background()
	suite.mockDynamoDb = mocks.NewDynamoDbProvider(suite.T())

	suite.lockManagerOptions = lock.LockManagerOptions{
		WithEncryption:   false,
		PartitionKeyName: suite.defaultPartitionKeyName,
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

	if suite.useEncryptionLockManager {
		suite.lockManagerOptions.WithEncryption = true
		options := lock.LockManagerOptionsWithEncryption{
			LockManagerOptions: suite.lockManagerOptions,
			EncryptionKey:      suite.encryptionKeyString,
		}
		suite.lockManager = lock.NewLockManagerWithEncryption(suite.mockDynamoDb, options)
	} else {
		suite.lockManager = lock.NewLockManager(suite.mockDynamoDb, suite.lockManagerOptions)
	}
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

func (suite *lock_manager_UnitTestSuite) Test_CanAcquireLock_NotOwner_Fail() {
	suite.setupMocks()
	suite.mockDynamoDb.EXPECT().
		Query(suite.mockContext, mock.Anything).
		Return(&dynamodb.QueryOutput{Count: 1, Items: suite.getAttributeValueMap("not_same_owner")}, nil).
		Once()

	err := suite.lockManager.CanAcquireLock(suite.mockContext, suite.defaultPartitionKey)
	suite.Error(err)
	suite.ErrorContains(err, "lock is in use by another and not released")
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

func (suite *lock_manager_UnitTestSuite) Test_AcquireLockWithEncryption_Success() {
	suite.useEncryptionLockManager = true
	suite.setupMocks()
	suite.mockDynamoDb.EXPECT().
		PutItem(suite.mockContext, mock.Anything).
		Return(&dynamodb.PutItemOutput{Attributes: suite.getAttributeValueMap(suite.defaultLockOwner)[0]}, nil).
		Once()

	options := &models.AcquireLockOptions{
		PartitionKey: suite.defaultPartitionKey,
		Data:         suite.dataOption,
	}

	suite.useEncryptionLockManager = false

	lk, err := suite.lockManager.AcquireLock(suite.mockContext, options)
	suite.Nil(err)
	suite.NotNil(lk)
}

func (suite *lock_manager_UnitTestSuite) Test_ReleaseLock_NotDeleteOnRelease_Success() {
	suite.setupMocks()
	output := &dynamodb.GetItemOutput{
		Item: suite.getAttributeValueMap(suite.defaultLockOwner)[0],
	}
	suite.mockDynamoDb.EXPECT().
		GetItem(suite.mockContext, mock.Anything).
		Return(output, nil).
		Once()
	suite.mockDynamoDb.EXPECT().
		UpdateItem(suite.mockContext, mock.Anything).
		Return(&dynamodb.UpdateItemOutput{}, nil).
		Once()

	result, err := suite.lockManager.ReleaseLock(suite.mockContext, suite.defaultLockOwner)
	suite.Nil(err)
	suite.True(result)
}

func (suite *lock_manager_UnitTestSuite) Test_ReleaseLock_DeleteOnRelease_Success() {
	suite.deleteOnReleaseOption = true
	suite.setupMocks()
	output := &dynamodb.GetItemOutput{
		Item: suite.getAttributeValueMap(suite.defaultLockOwner)[0],
	}
	suite.mockDynamoDb.EXPECT().
		GetItem(suite.mockContext, mock.Anything).
		Return(output, nil).
		Once()
	suite.mockDynamoDb.EXPECT().
		DeleteItem(suite.mockContext, mock.Anything).
		Return(&dynamodb.DeleteItemOutput{}, nil).
		Once()
	suite.deleteOnReleaseOption = false

	result, err := suite.lockManager.ReleaseLock(suite.mockContext, suite.defaultLockOwner)
	suite.Nil(err)
	suite.True(result)
}

func (suite *lock_manager_UnitTestSuite) Test_CreateTable() {
	suite.setupMocks()
	suite.mockDynamoDb.EXPECT().
		CreateTable(suite.mockContext, mock.Anything).
		Return(&dynamodb.CreateTableOutput{}, nil).
		Once()

	tableOptions := &models.CreateDynamoDBTableOptions{
		TableName:    suite.defaultTableName,
		PartitionKey: suite.defaultPartitionKeyName,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  nil,
			WriteCapacityUnits: nil,
		},
	}

	result, err := suite.lockManager.CreateLockTable(suite.mockContext, tableOptions)
	suite.Nil(err)
	suite.NotNil(result)
}

func (suite *lock_manager_UnitTestSuite) Test_GetLock_Success() {
	suite.setupMocks()
	suite.mockDynamoDb.EXPECT().
		GetItem(suite.mockContext, mock.Anything).
		Return(&dynamodb.GetItemOutput{Item: suite.getAttributeValueMap(suite.defaultLockOwner)[0]}, nil).
		Once()

	result, err := suite.lockManager.GetLock(suite.mockContext, suite.defaultLockOwner)

	log.Printf("%+v", result)
	suite.Nil(err)
	suite.NotNil(result)
}

func (suite *lock_manager_UnitTestSuite) Test_GetLockWithEncryption_Success() {
	suite.useEncryptionLockManager = true
	suite.setupMocks()
	suite.mockDynamoDb.EXPECT().
		GetItem(suite.mockContext, mock.Anything).
		Return(&dynamodb.GetItemOutput{Item: suite.getAttributeValueMap(suite.defaultLockOwner)[0]}, nil).
		Once()
	suite.useEncryptionLockManager = false

	result, err := suite.lockManager.GetLock(suite.mockContext, suite.defaultLockOwner)

	log.Printf("%+v", result)
	suite.Nil(err)
	suite.NotNil(result)
}

func (suite *lock_manager_UnitTestSuite) getAttributeValueMap(ownerName string) []map[string]types.AttributeValue {
	return []map[string]types.AttributeValue{
		{
			suite.defaultPartitionKeyName: &types.AttributeValueMemberS{Value: suite.defaultPartitionKey},
			"owner":                       &types.AttributeValueMemberS{Value: ownerName},
			"timestamp":                   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", suite.lockOutputValue.Timestamp)},
			"ttl":                         &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", suite.lockOutputValue.Ttl)},
			"isReleased":                  &types.AttributeValueMemberBOOL{Value: suite.isReleasedOption},
			"deleteOnRelease":             &types.AttributeValueMemberBOOL{Value: suite.deleteOnReleaseOption},
			"data":                        &types.AttributeValueMemberS{Value: string(suite.dataOption)},
		},
	}
}
