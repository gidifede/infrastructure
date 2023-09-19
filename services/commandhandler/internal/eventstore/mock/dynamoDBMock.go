package mockeventstore

import (
	"command-handler/internal/utils"
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type MockDynamoDB struct {
	dynamodbiface.DynamoDBAPI
}

func (m *MockDynamoDB) QueryWithContext(ctx context.Context, input *dynamodb.QueryInput, options ...request.Option) (*dynamodb.QueryOutput, error) {

	fmt.Printf("ctx :%v , request: %v",ctx,options)
	var items []map[string]*dynamodb.AttributeValue
	var output dynamodb.QueryOutput
	requestedAggregateID := *(input.ExpressionAttributeValues[":aggregate_id"]).S

	if requestedAggregateID == "idNotFound" {
		count := int64(0)
		output = dynamodb.QueryOutput{
			Items: items,
			Count: &count,
		}
	} else {
		items = append(items, buildFakeItem("C234-1234-1248", 1673885460, "Logistic.PCL.Product.Accept.AcceptProduct", "A json", "Logistic.PCL.UP.OMP", "eventId1", 1673885460))
		items = append(items, buildFakeItem("C234-1234-1248", 1673996280, "Logistic.PCL.Product.Withdraw.WithdrawProduct", "Another json", "Logistic.PCL.UP.ANOTHER", "eventId2", 1673996280))
		count := int64(2)
		output = dynamodb.QueryOutput{
			Items: items,
			Count: &count,
		}
	}
	return &output, nil
}

func (m *MockDynamoDB) PutItemWithContext(ctx context.Context, input *dynamodb.PutItemInput, options ...request.Option) (*dynamodb.PutItemOutput, error) {
	fmt.Printf("ctx :%v , request: %v",ctx,options)
	if *(input.Item["aggregate_id"]).S == "aggregateOK" {
		return nil, nil
	} 

	return nil, fmt.Errorf("an error")
}

func buildFakeItem(aggregateID string, timestamp int64, aggregateType string, data string, source string, eventID string, timestampSent int64) map[string]*dynamodb.AttributeValue {
	item := make(map[string]*dynamodb.AttributeValue)
	aggregateIDAttribute := dynamodb.AttributeValue{}
	aggregateIDAttribute.SetS(aggregateID)
	timestampAttribute := dynamodb.AttributeValue{}
	timestampAttribute.SetN(strconv.FormatInt(timestamp, 10))
	typeAttribute := dynamodb.AttributeValue{}
	typeAttribute.SetS(aggregateType)
	dataAttribute := dynamodb.AttributeValue{}
	dataAttribute.SetS(data)
	sourceAttribute := dynamodb.AttributeValue{}
	sourceAttribute.SetS(source)
	eventIDAttribute := dynamodb.AttributeValue{}
	eventIDAttribute.SetS(eventID)
	timestampSentAttribute := dynamodb.AttributeValue{}
	timestampSentAttribute.SetN(strconv.FormatInt(timestampSent, 10))
	dataChecksumAttribute := dynamodb.AttributeValue{}
	dataChecksumAttribute.SetS(utils.ComputeChecksum([]byte(data), true))
	item["aggregate_id"] = &aggregateIDAttribute
	item["timestamp"] = &timestampAttribute
	item["type"] = &typeAttribute
	item["data"] = &dataAttribute
	item["source"] = &sourceAttribute
	item["eventId"] = &eventIDAttribute
	item["timestampSent"] = &timestampSentAttribute
	item["dataChecksum"] = &dataChecksumAttribute
	return item
}
