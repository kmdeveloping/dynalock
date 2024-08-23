package lock

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kmdeveloping/dynalock/internal/models"
	"log"
	"strconv"
)

const (
	attrData            = "data"
	attrOwnerName       = "owner"
	attrIsReleased      = "isReleased"
	attrTimestamp       = "timestamp"
	attrExpirationTime  = "ttl"
	attrDeleteOnRelease = "deleteOnRelease"
)

func (lm *LockManager) marshalLockItem(item models.Lock) (map[string]types.AttributeValue, error) {
	av := map[string]types.AttributeValue{
		lm.partitionKeyName: &types.AttributeValueMemberS{Value: item.PartitionKey},
		attrOwnerName:       &types.AttributeValueMemberS{Value: item.Owner},
		attrTimestamp:       &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", item.Timestamp)},
		attrExpirationTime:  &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", item.Ttl)},
		attrDeleteOnRelease: &types.AttributeValueMemberBOOL{Value: item.DeleteOnRelease},
		attrIsReleased:      &types.AttributeValueMemberBOOL{Value: item.IsReleased},
		attrData:            &types.AttributeValueMemberS{Value: string(item.Data)},
	}

	return av, nil
}

func (lm *LockManager) unmarshalLockItem(item map[string]types.AttributeValue) (models.Lock, error) {
	timestamp, _ := strconv.ParseInt(item[attrTimestamp].(*types.AttributeValueMemberN).Value, 10, 64)
	expTime, _ := strconv.ParseInt(item[attrExpirationTime].(*types.AttributeValueMemberN).Value, 10, 64)

	result := models.Lock{
		PartitionKey:    item[lm.partitionKeyName].(*types.AttributeValueMemberS).Value,
		Owner:           item[attrOwnerName].(*types.AttributeValueMemberS).Value,
		DeleteOnRelease: item[attrDeleteOnRelease].(*types.AttributeValueMemberBOOL).Value,
		IsReleased:      item[attrIsReleased].(*types.AttributeValueMemberBOOL).Value,
		Data:            []byte(item[attrData].(*types.AttributeValueMemberS).Value),
		Timestamp:       timestamp,
		Ttl:             expTime,
	}

	return result, nil
}

func (lm *LockManager) unmarshalListLocks(items []map[string]types.AttributeValue) ([]models.Lock, error) {
	var result []models.Lock

	for _, item := range items {
		lock, err := lm.unmarshalLockItem(item)
		if err != nil {
			log.Printf("failed to unmarshal lock item: %v", err)
			continue
		}

		result = append(result, lock)
	}

	return result, nil
}

func (lm *LockManager) isDeleteOnRelease(item map[string]types.AttributeValue) (bool, error) {
	return item[attrDeleteOnRelease].(*types.AttributeValueMemberBOOL).Value, nil
}
