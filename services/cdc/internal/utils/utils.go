package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	overlog "github.com/Trendyol/overlog"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-xray-sdk-go/xray"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog/log"
)

func AddClassAndMethodToMDC(i interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	functionName := runtime.FuncForPC(pc).Name()

	str := strings.Split(functionName, ".")
	method := str[len(str)-1]
	// class := str[len(str)-2]
	class := reflect.Indirect(reflect.ValueOf(i)).Type().Name()

	overlog.MDC().Set("method", method)
	overlog.MDC().Set("class", class)

}

func ConvertToCloudEvent(body map[string]string) (*cloudevents.Event, error) {
	event := cloudevents.NewEvent()

	event.SetID(body["eventId"])
	data, err := strconv.Unquote(body["data"])
	if nil != err {
		// invalid sintax error, nothing to unquote
		data = body["data"]
	}
	dataStructured := map[string]interface{}{}
	err = json.Unmarshal([]byte(data), &dataStructured)
	if nil != err {
		return nil, err
	}
	event.SetData(cloudevents.ApplicationJSON, dataStructured)
	event.SetSource(body["source"])
	event.SetType(body["type"])

	i, err := strconv.ParseInt(body["timestampSent"], 10, 64)
	if nil != err {
		return nil, err
	}
	t := time.UnixMilli(i)

	event.SetTime(t)
	return &event, nil
}

func ParseDynamoStream(ctx context.Context, e events.DynamoDBEventRecord, jsonm map[string]string) (string, error) {

	xray.BeginSubsegment(ctx, "Dynamo Parsing...")
	// jsonm := make(map[string]string)

	// for _, record := range e.Records {
	log.Debug().Msg(fmt.Sprintf("Processing request data for event ID %s, type %s.\n", e.EventID, e.EventName))

	for name, value := range e.Change.NewImage {

		if value.DataType() == events.DataTypeNumber {
			jsonm[name] = value.Number()
		} else {
			jsonm[name] = value.String()
		}
	}
	// }
	msgByte, err := json.Marshal(jsonm)
	if err != nil {
		log.Err(err)
		return "", err
	}

	msgJSON := string(msgByte)

	log.Debug().Msg(fmt.Sprintf("Message ready to be sent: %s", msgJSON))

	return msgJSON, nil
}
