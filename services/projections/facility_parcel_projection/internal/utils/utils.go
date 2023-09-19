package utils

import (
	"reflect"
	"runtime"
	"strings"

	overlog "github.com/Trendyol/overlog"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-xray-sdk-go/header"
)

func GetXRayTraceID(event events.SQSMessage) string {
	traceID := header.FromString(event.Attributes["AWSTraceHeader"]).TraceID
	return traceID
}

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
