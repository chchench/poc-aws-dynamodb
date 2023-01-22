package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	AWS_REGION = "us-east-1"

	TABLE_NAME               = "RequestRecords"
	ITEM_TO_WRITE_ID         = "1234567"
	ITEM_TO_WRITE_COMPANY_ID = "0987654321"
)

func main() {
	client := setupClient()

	putItem(client, ITEM_TO_WRITE_ID, "2023-01-22T09:52:23.616414-08:00")
	putItem(client, ITEM_TO_WRITE_ID, "")

	records, _ := queryItems(client, ITEM_TO_WRITE_ID)
	for i, rec := range records {
		fmt.Printf("Record #%d: %v\n", i, rec)
	}

	before := len(records)
	fmt.Printf("# records BEFORE deletion = %d\n", before)

	deleteItem(client, ITEM_TO_WRITE_ID, "2023-01-22T09:52:23.616414-08:00")

	records, _ = queryItems(client, ITEM_TO_WRITE_ID)
	for i, rec := range records {
		fmt.Printf("Record #%d: %v\n", i, rec)
	}

	after := len(records)
	fmt.Printf("# records AFTER deletion = %d\n", after)

	if before != (after + 1) {
		fatalExit("after deletion record count is not correct: expected %d, got %d\n",
			before-1, after)
	}

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

func putItem(client *dynamodb.Client, itemId string, timestamp string) {

	if timestamp == "" {
		currentTime := time.Now()
		timestamp = currentTime.Format(time.RFC3339Nano)
	}

	_, err := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),

		Item: map[string]types.AttributeValue{
			"Id":        &types.AttributeValueMemberS{Value: itemId},
			"Timestamp": &types.AttributeValueMemberS{Value: timestamp},
			"CompanyId": &types.AttributeValueMemberS{Value: ITEM_TO_WRITE_COMPANY_ID},
			"JSON":      &types.AttributeValueMemberS{Value: "{ \"key\": \"value-" + timestamp + "\" }"},
		},
	})

	if err != nil {
		fatalExit("encountered problem while putting item to table: %v\n", err)
	}
}

/* Attempt to delete an item from the table. The key thing here is to make sure
the filtering criteria match the schema. */

func deleteItem(client *dynamodb.Client, itemId string, timestamp string) {

	_, err := client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(TABLE_NAME),

		Key: map[string]types.AttributeValue{
			"Id":        &types.AttributeValueMemberS{Value: itemId},
			"Timestamp": &types.AttributeValueMemberS{Value: timestamp},
		},
	})
	if err != nil {
		fatalExit("encountered problem while deleting item from table: %v\n", err)
	}
}

/* Query and retrieve one or more items with the same "Id" */

type RequestRecord struct {
	Id        string `dynamodbav:"id"`
	Timestamp string `dynamodbav:"timestamp"`
	JSON      string `dynamodbav:"json"`
}

func queryItems(client *dynamodb.Client, id string) ([]RequestRecord, error) {

	var err error
	var response *dynamodb.QueryOutput

	var records []RequestRecord

	keyEx := expression.Key("Id").Equal(expression.Value(id))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()

	if err != nil {
		log.Printf("Couldn't build epxression for query. Here's why: %v\n", err)
	} else {
		response, err = client.Query(context.TODO(), &dynamodb.QueryInput{
			TableName:                 aws.String(TABLE_NAME),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			KeyConditionExpression:    expr.KeyCondition(),
		})
		if err != nil {
			log.Printf("unable to query for records based on Id: %v: ", err)
		} else {
			err = attributevalue.UnmarshalListOfMaps(response.Items, &records)
			if err != nil {
				log.Printf("unable to unmarshal query response: %v\n", err)
			}
		}
	}
	return records, err
}

func fatalExit(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
