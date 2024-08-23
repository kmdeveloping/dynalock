package dynalock_tests

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kmdeveloping/dynalock"
	"github.com/kmdeveloping/dynalock/internal/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type dynalock_client_UnitTestSuite struct {
	suite.Suite

	mockDynamoDb *mocks.DynamoDbProvider
	mockContext  context.Context

	client           *dynalock.Client
	defaultTableName string
	defaultOwnerName string
}

func Test_dynalock_client_UnitTestSuite(t *testing.T) {
	suite.Run(t, &dynalock_client_UnitTestSuite{})
}

func (suite *dynalock_client_UnitTestSuite) SetupTest() {
	suite.defaultTableName = "dynalock_default_table"
	suite.defaultOwnerName = "dynalock_default_owner"
}

func (suite *dynalock_client_UnitTestSuite) setupMocks() {
	suite.mockContext = context.Background()
	suite.mockDynamoDb = mocks.NewDynamoDbProvider(suite.T())

	suite.client = dynalock.NewDynalockClient(suite.mockDynamoDb,
		suite.defaultTableName,
		dynalock.WithLockOwnerName(suite.defaultOwnerName))
}

func (suite *dynalock_client_UnitTestSuite) Test_AcquireLock_Success() {
	suite.setupMocks()
	suite.mockDynamoDb.EXPECT().
		Query(suite.mockContext, mock.Anything).
		Return(&dynamodb.QueryOutput{Items: nil}, nil).
		Once()
	suite.mockDynamoDb.EXPECT().
		PutItem(suite.mockContext, mock.Anything).
		Return(nil, nil).
		Once()

	lock, err := suite.client.AcquireLock(suite.mockContext, "text_key_action")
	suite.NoError(err)
	suite.NotNil(lock)
}

func (suite *dynalock_client_UnitTestSuite) Test_AcquireLock_NotReleased_Fail() {
	key := "mock:key_data"
	items := []map[string]types.AttributeValue{
		{
			"key":             &types.AttributeValueMemberS{Value: key},
			"owner":           &types.AttributeValueMemberS{Value: "ima_different_owner"},
			"timestamp":       &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", time.Now().Unix())},
			"ttl":             &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", 30*24*60*60)},
			"isReleased":      &types.AttributeValueMemberBOOL{Value: false},
			"deleteOnRelease": &types.AttributeValueMemberBOOL{Value: false},
			"data":            &types.AttributeValueMemberS{Value: `{"key":"value"}`},
		},
	}

	suite.setupMocks()
	suite.mockDynamoDb.EXPECT().
		Query(suite.mockContext, mock.Anything).
		Return(&dynamodb.QueryOutput{Count: 1, Items: items}, nil).
		Once()

	lock, err := suite.client.AcquireLock(suite.mockContext, key)
	suite.Nil(lock)
	suite.Error(err)
}
