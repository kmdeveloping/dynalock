package lock

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kmdeveloping/dynalock/internal/models"
)

func marshalLockItem(item models.Lock) (map[string]types.AttributeValue, error) {
	av := map[string]types.AttributeValue{
		DefaultPartitionKeyName: &types.AttributeValueMemberS{Value: item.PartitionKey},
		attrOwnerName:           &types.AttributeValueMemberS{Value: item.Owner},
		attrTimestamp:           &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", item.Timestamp)},
		attrExpirationTime:      &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", item.Ttl)},
		attrDeleteOnRelease:     &types.AttributeValueMemberBOOL{Value: item.DeleteOnRelease},
		attrIsReleased:          &types.AttributeValueMemberBOOL{Value: item.IsReleased},
		attrData:                &types.AttributeValueMemberS{Value: string(item.Data)},
	}

	return av, nil
}

func unmarshalLockItem(item map[string]types.AttributeValue) (models.Lock, error) {

}

func unmarshalListLocks(items []map[string]types.AttributeValue) ([]models.Lock, error) {

}
