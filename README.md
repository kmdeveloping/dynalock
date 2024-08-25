# dynalock

### Installation:
* ```go get github.com/kmdeveloping/dynalock```

### Usage:

```go
package main

import (
	"context"
	"log"

	// all imports for your app
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/kmdeveloping/dynalock"
)

func main() {
	// this is your aws connection config, use your relevant config setup
	cfg := aws.NewConfig()
	
	// create a standard client for dynamodb using config settings
	dynamoClient := dynamodb.NewFromConfig(cfg)
	
	// create the lock client 
	// options like dynalock.WithLockOwnerName are optional as omitting will generate random owner name
	// dynalock.WithPartitionKeyName will override the default key name of "key"
	client := dynalock.NewDynalockClient(dynamoClient,
		"default_table_name_of_your_choice",
		dynalock.WithLockOwnerName("optionally_set_your_custom_lock_owner_name"),
		dynalock.WithPartitionKeyName("optionally_overwrite_partition_key_name_default"))

	// metadata is optional to store encrypted OR unencrypted within the lock for using later
	// use the dynalock.WithData option to pass in a []byte() data piece
	someMetaData := []byte(`{"meta": "data"}`)
	
	// acquire the lock with a key name and add data (optional to add data)
	// dynalock.DeleteOnRelease is what it sounds like, it will clear the lock from the lock table on release 
	// use this to keep tables clean if needed, understand this WILL DELETE the lock with data from the lock table
	acquiredLock, err := client.AcquireLock(context.Background(),
		"some_key_name",
		dynalock.WithData(someMetaData),
		dynalock.DeleteOnRelease())

	if err != nil {
		// do something with the error
		log.Printf("%v", err)
		return
	}
	
	log.Printf("%+v", acquiredLock)
	
	// at this point the lock is acquired and will need to be released to be actioned by another caller.
	// because the lock is stored in dynamo, it is not affected by lambda's runtime ending. this means another 
	// lambda (using the same owner name) can update or release the lock as needed making this a true distributed locking client
	
	// release the lock using the key name of the lock 
	released, err := client.ReleaseLock(context.Background(), "some_key_name")
	if err != nil {
		// do something with the error
		log.Printf("%v", err)
		return
    }
	
	// if true, the lock is released for another client to acquire, 
	// the lock is only released if the owner matches the client owner 
	// ReleaseLock will return error is the owner is mismatched
	log.Println(released)
}

```