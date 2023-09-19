package validator

import (
	"io/ioutil"
	"os"
	"testing"

	overlog "github.com/Trendyol/overlog"
	"github.com/rs/zerolog"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zlogger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	overlog.New(zlogger)
	overlog.SetGlobalFields([]string{"trace id", "method", "class"})

	overlog.MDC().Set("trace id", "0123456789")
}

func readRefFile() []byte {
	jsonFilePath := "testData/common.json"

	jsonFile, errOpenFile := os.Open(jsonFilePath)
	if errOpenFile != nil {
		panic("Error in retrieving json Schema file")
	}

	defer jsonFile.Close()

	byteValue, errReadFile := ioutil.ReadAll(jsonFile)
	if errReadFile != nil {
		panic("Error in reading json Schema file")
	}

	return byteValue
}

func readFile() []byte {

	jsonFilePath := "testData/JSONSchema.json"

	jsonFile, errOpenFile := os.Open(jsonFilePath)
	if errOpenFile != nil {
		panic("Error in retrieving json Schema file")
	}

	defer jsonFile.Close()

	byteValue, errReadFile := ioutil.ReadAll(jsonFile)
	if errReadFile != nil {
		panic("Error in reading json Schema file")
	}

	return byteValue

}

func TestValidate(t *testing.T) {

	type args struct {
		body           []byte
		bodyJSONSchema []byte
	}

	tests := []struct {
		name          string
		args          args
		wantErr       bool
		wantErrString string
	}{
		{"Validate Command Accept OK", args{[]byte(`{
			"product":{
			   "name": "string",
			   "id": "string",
			   "type": "BOX",
			   "attributes": []
		   },
			 "location" : {
				 "type" : "string",
				 "address" : "string",
				 "zipcode" : "00100",
				 "city" : "string",
				 "nation" : "string",
			   "locationCode": "string",
			   "attributes": []
			 },
			 "sender":{
				"name":"string",
				"province":"string",
				"city":"string",
				"address":"string",
				"zipcode":"00100",
				"attributes": []
			 },
			 "receiver":{
				"name":"string",
				"province":"string",
				"city":"string",
				"address":"string",
				"zipcode":"00100",
				"number":"string",
				"email":"test@poste.it",
				"note":"string",
				"attributes": []
			},
			"timestamp" : "2018-11-13T20:20:39+00:00",
			"attributes": []
		 }
		 `), readFile()}, false, ""},
		{"Validate Command Accept KO", args{[]byte(`{
					"product":{
					   "name": "string",
					   "id": "string",
					   "type": "string",
					   "attributes": []
				   },
					 "location" : {
						 "type" : "string",
						 "address" : "string",
						 "zipcode" : "string",
						 "city" : "string",
						 "nation" : "string",
					   "locationCode": "string",
					   "attributes": []
					 },
					 "sender":{
						"name":"string",
						"province":"string",
						"city":"string",
						"address":"string",
						"zipcode":"string",
						"attributes": []
					 },
					"timestamp" : "2018-11-13T20:20:39+00:00",
					"attributes": []
				 }
				 `), readFile()}, true, "validation error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Validator{}
			err := v.ValidateCommand(tt.args.body, tt.args.bodyJSONSchema, readRefFile())
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err != nil) && err.Error() != tt.wantErrString {
				t.Errorf("ValidateCommand() errorString = %v, wantErrString = %v", err.Error(), tt.wantErrString)
			}
		})
	}

}
