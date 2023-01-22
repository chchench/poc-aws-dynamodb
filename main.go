package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

const (
	AWS_REGION = "us-east-1"

	TABLE_NAME                       = "RequestRecords"
	EXAMPLE_ITEM_TO_WRITE_ID         = "1234567"
	EXAMPLE_ITEM_TO_WRITE_COMPANY_ID = "0987654321"
)

func main() {
	client := setupClient()

	putItem(client, "2023-01-22T09:52:23.616414-08:00")
	putItem(client, "2023-01-22T09:52:23.616414-08:00")
	putItem(client, "")
	putItem(client, "")

	deleteItem(client, "2023-01-22T09:52:23.616414-08:00")
}

func setupClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = AWS_REGION
		return nil
	})
	if err != nil {
		fatalExit("encountered problem while configuring client: %v\n", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return client
}

/* Attempt to insert an item into the table. In our case, the primary key used to
determine uniqueness of an entry, setup for the table using the terraform script, is
"Id" + "Timestamp" so if these two values are same as those for an existing item in
the table, then the subsequent item won't be added to the table. */

func putItem(client *dynamodb.Client, timestamp string) {

	if timestamp == "" {
		currentTime := time.Now()
		timestamp = currentTime.Format(time.RFC3339Nano)
	}

	output, err := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),

		Item: map[string]types.AttributeValue{
			"Id":        &types.AttributeValueMemberS{Value: EXAMPLE_ITEM_TO_WRITE_ID},
			"Timestamp": &types.AttributeValueMemberS{Value: timestamp},
			"CompanyId": &types.AttributeValueMemberS{Value: EXAMPLE_ITEM_TO_WRITE_COMPANY_ID},
			"JSON":      &types.AttributeValueMemberS{Value: "{ \"test\": \"testvalue" + timestamp + "\" }"},
		},
	})

	if err != nil {
		fatalExit("encountered problem while putting item to table: %v\n", err)
	}

	fmt.Println(output.Attributes)
}

/* Attempt to delete an item from the table. The key thing here is to make sure
the filtering criteria match the schema. */

func deleteItem(client *dynamodb.Client, timestamp string) {

	output, err := client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(TABLE_NAME),

		Key: map[string]types.AttributeValue{
			"Id":        &types.AttributeValueMemberS{Value: EXAMPLE_ITEM_TO_WRITE_ID},
			"Timestamp": &types.AttributeValueMemberS{Value: timestamp},
		},
	})
	if err != nil {
		fatalExit("encountered problem while deleting item from table: %v\n", err)
	}

	fmt.Println(output.Attributes)
}

func fatalExit(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
