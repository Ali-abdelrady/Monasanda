package convert

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"

	"mosanda-vacation-attendace/pkg/osioutil"

	"github.com/fatih/structs"
)

func MapToObject(parSource map[string]interface{}, parDestination any) error {

	pvDestinationValue := reflect.ValueOf(parDestination)
	if pvDestinationValue.Kind() != reflect.Ptr {
		return osioutil.GetRefErrorTxt("MPTST0001", "destination must be pointer")
	}

	pvDestinationElem := pvDestinationValue.Elem()

	if pvDestinationElem.Kind() == reflect.Map {
		pvType := pvDestinationElem.Type()

		fmt.Printf("UUUUUUUUUUUUUUUUUUUU:02, [%v]\n", pvType)

		var pvError error

		pvDestination, pvError := mapToMap(parSource, pvDestinationValue, pvType)
		if pvError != nil {
			return pvError
		}

		parDestination = &pvDestination

		fmt.Printf("UUUUUUUUUUUUUUUUUUUU:03, [%T]\n", parDestination)

		return nil
	} else if pvDestinationElem.Kind() == reflect.Struct {
		return mapToStruct(parSource, parDestination)
	}

	return osioutil.GetRefErrorTxt("MPTST0002", "destination type must be either Map or Struct")
}

func mapToStruct(parSource map[string]interface{}, parDestination any) error {

	return nil
}

func mapToMap(parSource map[string]interface{}, parDestination reflect.Value, parDestType reflect.Type) (any, error) {
	// pmDestination := reflect.ValueOf(parDestination).Elem()
	// fmt.Println("YYYYYYYYYYYYYYYY:01", parDestination.(type))

	switch parDestType.String() {
	case "map[string]interface {}":
		// parDestination = &parSource
		return &parSource, nil

	default:
		// fmt.Printf("UUUUUUUUUUUUUUUUUUUU:01, [%T]\n", pmValue)
	}

	return nil, nil
}

func NewMapFromHttpRequest(parRequest *http.Request) (map[string]interface{}, error) {

	pvDataJSONBytes := []byte{}
	var pvError error

	if parRequest.Header.Get("Content-Encoding") == "gzip" {
		var pvReader *gzip.Reader
		pvReader, pvError = gzip.NewReader(parRequest.Body)
		if pvError != nil {
			return nil, osioutil.GetRefError("NMFHR0001", pvError)
		}

		if pvReader != nil {
			defer pvReader.Close()
		}

		if pvError == nil {
			pvDataJSONBytes, pvError = io.ReadAll(pvReader)
			if pvError != nil {
				return nil, osioutil.GetRefError("NMFHR0002", pvError)
			}
		} else {
			return nil, osioutil.GetRefError("NMFHR0003", pvError)
		}

	} else {
		pvDataJSONBytes, pvError = io.ReadAll(parRequest.Body)
		if pvError != nil {
			return nil, osioutil.GetRefError("NMFHR0004", pvError)
		}
	}

	parRequest.Body.Close()

	pmReturn := make(map[string]interface{})

	pvError = json.Unmarshal(pvDataJSONBytes, &pmReturn)
	if pvError != nil {
		return nil, osioutil.GetRefError("NMFHR0005", pvError)
	}

	return pmReturn, nil
}

func NewMapFromHttpResponse(parResponse *http.Response) (map[string]interface{}, error) {

	pvDataJSONBytes := []byte{}
	var pvError error

	if parResponse.Header.Get("Content-Encoding") == "gzip" {
		var pvReader *gzip.Reader
		pvReader, pvError = gzip.NewReader(parResponse.Body)
		if pvError == nil {
			pvDataJSONBytes, pvError = io.ReadAll(pvReader)
			if pvError != nil {
				return nil, osioutil.GetRefError("NMFHRS0001", pvError)
			}
		} else {
			return nil, osioutil.GetRefError("NMFHRS0002", pvError)
		}

		pvReader.Close()
	} else {
		pvDataJSONBytes, pvError = io.ReadAll(parResponse.Body)
		if pvError != nil {
			return nil, osioutil.GetRefError("NMFHRS0003", pvError)
		}
	}

	parResponse.Body.Close()

	psReturn := make(map[string]interface{})

	pvError = json.Unmarshal(pvDataJSONBytes, &psReturn)
	if pvError != nil {
		return nil, osioutil.GetRefError("NMFHRS0004", pvError)
	}

	return psReturn, nil
}

func NewMapFromJSONString(parJSON string) (map[string]interface{}, error) {
	pmReturn := make(map[string]interface{})

	pvError := json.Unmarshal([]byte(parJSON), &pmReturn)
	if pvError != nil {
		return nil, osioutil.GetRefError("NOFJS0001", pvError)
	}

	return pmReturn, nil
}

func NewMapFromJSONFile(parFilePath string) (map[string]interface{}, error) {
	pmReturn := make(map[string]interface{})

	b, pvError := os.ReadFile(parFilePath)
	if pvError != nil {
		return nil, osioutil.GetRefError("NODJF0001", pvError)
	}

	pvError = json.Unmarshal(b, &pmReturn)
	if pvError != nil {
		return nil, osioutil.GetRefError("NODJF0002", pvError)
	}

	return pmReturn, nil
}

func NewMapFromJSONBytes(parBytes []byte) (map[string]interface{}, error) {
	pmReturn := make(map[string]interface{})

	pvError := json.Unmarshal(parBytes, &pmReturn)
	if pvError != nil {
		return nil, osioutil.GetRefError("NOFJB0001", pvError)
	}

	return pmReturn, nil
}

func NewMapFromStruct(parStruct interface{}) map[string]interface{} {
	return structs.Map(parStruct)
}

func NewMapFromInterface(parInterface interface{}, parShare bool) (map[string]interface{}, error) {
	pmReturn := make(map[string]interface{})

	var pvInterface interface{}

	if fmt.Sprintf("%T", parInterface) == "map[string]interface {}" {
		pvInterface = parInterface
	} else {
		switch fmt.Sprintf("%T", parInterface) {
		case "string":
			pvError := json.Unmarshal([]byte(fmt.Sprint(parInterface)), pvInterface)
			if pvError != nil {
				return pmReturn, pvError
			}

		default:
			return pmReturn, nil
		}
	}

	if parShare {
		pmReturn = pvInterface.(map[string]interface{})
		return pmReturn, nil
	}

	b, pvError := json.Marshal(pvInterface)
	if pvError != nil {
		return pmReturn, pvError
	}

	pvError = json.Unmarshal(b, &pmReturn)
	if pvError != nil {
		return pmReturn, pvError
	}

	return pmReturn, nil
}
