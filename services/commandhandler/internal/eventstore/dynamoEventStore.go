package eventstore

import (
	utils "command-handler/internal/utils"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/rs/zerolog/log"
)

type DynamoEventStore struct {
	connection dynamodbiface.DynamoDBAPI
	tableName  string
}

func NewDynamoEventStore(conn dynamodbiface.DynamoDBAPI, tableName string) IEventStore {
	return &DynamoEventStore{connection: conn, tableName: tableName}
}

func (es *DynamoEventStore) Append(ctx context.Context, event Event) error {
	utils.AddClassAndMethodToMDC(es)

	av, err := dynamodbattribute.MarshalMap(event)
	if err != nil {
		return fmt.Errorf("got error marshalling item: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(es.tableName),
	}

	_, err = es.connection.PutItemWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("got error calling PutItem: %s", err)
	}

	log.Debug().Msgf("successfully added event with ID: '" + event.AggregateID + " to table: " + es.tableName)
	return nil
}

func (es *DynamoEventStore) Get(ctx context.Context, aggregateID string) ([]Event, error) {
	utils.AddClassAndMethodToMDC(es)
	result, err := es.connection.QueryWithContext(ctx, &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":aggregate_id": {
				S: aws.String(aggregateID),
			},
		},
		KeyConditionExpression: aws.String("aggregate_id = :aggregate_id"),
		TableName:              aws.String(es.tableName),
	})
	if err != nil {
		return nil, err
	}

	var items []Event
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal aggregate from event store. %v", err)
	}

	return items, nil
}
