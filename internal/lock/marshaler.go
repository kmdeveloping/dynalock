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

func (lm *LockManager) marshalLockItem(item models.Lock, output *map[string]types.AttributeValue) error {
	*output = map[string]types.AttributeValue{
		lm.partitionKeyName: &types.AttributeValueMemberS{Value: item.PartitionKey},
		attrOwnerName:       &types.AttributeValueMemberS{Value: item.Owner},
		attrTimestamp:       &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", item.Timestamp)},
		attrExpirationTime:  &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", item.Ttl)},
		attrDeleteOnRelease: &types.AttributeValueMemberBOOL{Value: item.DeleteOnRelease},
		attrIsReleased:      &types.AttributeValueMemberBOOL{Value: item.IsReleased},
		attrData:            &types.AttributeValueMemberS{Value: string(item.Data)},
	}

	return nil
}

func (lm *LockManager) unmarshalLockItem(item map[string]types.AttributeValue, output *models.Lock) error {
	timestamp, _ := strconv.ParseInt(item[attrTimestamp].(*types.AttributeValueMemberN).Value, 10, 64)
	expTime, _ := strconv.ParseInt(item[attrExpirationTime].(*types.AttributeValueMemberN).Value, 10, 64)

	*output = models.Lock{
		PartitionKey:    item[lm.partitionKeyName].(*types.AttributeValueMemberS).Value,
		Owner:           item[attrOwnerName].(*types.AttributeValueMemberS).Value,
		DeleteOnRelease: item[attrDeleteOnRelease].(*types.AttributeValueMemberBOOL).Value,
		IsReleased:      item[attrIsReleased].(*types.AttributeValueMemberBOOL).Value,
		Data:            []byte(item[attrData].(*types.AttributeValueMemberS).Value),
		Timestamp:       timestamp,
		Ttl:             expTime,
	}

	return nil
}

func (lm *LockManager) unmarshalListLocks(items []map[string]types.AttributeValue, output *[]models.Lock) error {
	for _, item := range items {
		lock := models.Lock{}
		err := lm.unmarshalLockItem(item, &lock)
		if err != nil {
			log.Printf("failed to unmarshal lock item: %v", err)
			continue
		}

		*output = append(*output, lock)
	}

	return nil
}
