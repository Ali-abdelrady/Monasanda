package xhttputil

import (
	"bytes"
	"encoding/json"
	"time"

	// "mosanda-vacation-attendace/pkg/osioutil"

	"fmt"
	"net/http"

	"strconv"

	"mosanda-vacation-attendace/pkg/convert"
	"mosanda-vacation-attendace/pkg/osioutil"
)

var lvClient *http.Client

func SendResponse(parResponse http.ResponseWriter, parRequest *http.Request, parData map[string]interface{}) {
	pvPath := parRequest.URL.Path[1:]

	// b := parData.ExportAsJSONBytes()
	b, pvError := json.Marshal(parData)
	if pvError != nil {
		fmt.Println(pvError)
	}

	parResponse.Header().Set("Content-Type", "application/json")
	parResponse.Header().Set("Access-Control-Allow-Origin", "*")
	parResponse.Header().Set("Access-Control-Allow-Credentials:", "true")
	parResponse.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	parResponse.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Authorization, Origin, Content-Type, Version, content-encoding, apikey")
	parResponse.Header().Set("Access-Control-Expose-Headers", "X-Requested-With, Authorization, Origin, Content-Type, Version, content-encoding, apikey")
	parResponse.Header().Set("Access-Control-Max-Age", "3600")
	parResponse.Header().Set("Accept-Ranges", "bytes")
	parResponse.Header().Set("Content-Encoding", "gzip")

	b1, pvError := osioutil.GzipCompressData(b)
	if pvError != nil {
		fmt.Println(pvError)
	} else {
		parResponse.Header().Set("Content-Length", strconv.Itoa(len(b1)))
		i, pvError := parResponse.Write(b1)
		if pvError != nil {
			fmt.Println(pvError)
		}

		if i == len(b1) {
			// tlogger.PrintApiLog(pvPath)
			fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.999"), " API1 "+pvPath+".")
		} else {
			// tlogger.PrintApiLog(pvPath + " ERROR")
			fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.999"), " API2 "+pvPath+" ERROR.")
		}
	}
}

func HandleResponseError(parResponse http.ResponseWriter, parRequest *http.Request, parResponseData map[string]interface{}, parErrorRef string, parError error, parLang string, parURL string) {

	pvRequest := make(map[string]interface{})
	pvRequest["code"] = parErrorRef
	pvRequest["lang"] = parLang
	pvHeaders := make(map[string]string)
	pvHeaders["Content-Type"] = "application/json"
	pvResponse, _ := PostMapRequest(parURL+"encdec/getErrorMsg", pvRequest, pvHeaders)
	pvErrorMsg, _ := pvResponse["msg"]

	parResponseData["hasError"] = true
	parResponseData["ErrorToDisplay"] = pvErrorMsg
	parResponseData["ErrorMessage"] = fmt.Sprint(parError)
	parResponseData["ErrorRef"] = parErrorRef

	SendResponse(parResponse, parRequest, parResponseData)
}

func HandleResponseErrorText(parResponse http.ResponseWriter, parRequest *http.Request, parResponseData map[string]interface{}, parErrorRef string, parError string, parLang string, parURL string) {

	pvRequest := make(map[string]interface{})
	pvRequest["code"] = parErrorRef
	pvRequest["lang"] = parLang
	pvHeaders := make(map[string]string)
	pvHeaders["Content-Type"] = "application/json"
	pvResponse, _ := PostMapRequest(parURL+"encdec/getErrorMsg", pvRequest, pvHeaders)
	pvErrorMsg, _ := pvResponse["msg"]

	parResponseData["hasError"] = true
	parResponseData["ErrorToDisplay"] = pvErrorMsg
	parResponseData["ErrorMessage"] = parError
	parResponseData["ErrorRef"] = parErrorRef

	SendResponse(parResponse, parRequest, parResponseData)
}

func HandleResponseErrorTextMsg(parResponse http.ResponseWriter, parRequest *http.Request, parErrorRef string, parError string) {
	pmResponseData := make(map[string]interface{})

	pmResponseData["hasError"] = true
	pmResponseData["ErrorMessage"] = parError
	pmResponseData["ErrorRef"] = parErrorRef

	SendResponse(parResponse, parRequest, pmResponseData)
}

func HandleResponseErrorMsg(parResponse http.ResponseWriter, parRequest *http.Request, parErrorRef string, parError error) {
	pmResponseData := make(map[string]interface{})

	pmResponseData["hasError"] = true
	pmResponseData["ErrorMessage"] = fmt.Sprint(parError)
	pmResponseData["ErrorRef"] = parErrorRef

	SendResponse(parResponse, parRequest, pmResponseData)
}

func PostMapRequest(parURL string, parRequest map[string]interface{}, parHeaders map[string]string) (map[string]interface{}, error) {
	//https://hackernoon.com/avoiding-memory-leak-in-golang-api-1843ef45fca8

	b, pvError := json.Marshal(parRequest)
	if pvError != nil {
		return nil, osioutil.GetRefError("HTPU0001", pvError)
	}

	return PostBytesRequest(parURL, b, parHeaders)
}

func PostBytesRequest(parURL string, parRequest []byte, parHeaders map[string]string) (map[string]interface{}, error) {
	//https://hackernoon.com/avoiding-memory-leak-in-golang-api-1843ef45fca8

	if parHeaders == nil {
		parHeaders = make(map[string]string)
	}

	var pvBuff *bytes.Buffer

	if len(parRequest) == 0 {
		parRequest = []byte("{}")
	}

	b, _ := osioutil.GzipCompressData(parRequest)
	pvBuff = bytes.NewBuffer(b)

	pvRequest, pvError := http.NewRequest("POST", parURL, pvBuff)
	if pvError != nil {
		return map[string]interface{}{}, fmt.Errorf("[PBR0001] %v", pvError)
	}

	for pvKey, pvValue := range parHeaders {
		pvRequest.Header.Set(pvKey, pvValue)
	}
	pvRequest.Header.Set("Content-Encoding", "gzip")

	if lvClient == nil {

		pvDefaultRoundTripper := http.DefaultTransport
		pvDefaultTransportPointer, ok := pvDefaultRoundTripper.(*http.Transport)
		if !ok {
			return map[string]interface{}{}, fmt.Errorf("[PBR0002] DefaultRoundTripper not an *http.Transport")
		}

		pvDefaultTransport := *pvDefaultTransportPointer // dereference it to get a copy of the struct that the pointer points to
		pvDefaultTransport.MaxIdleConns = 3
		pvDefaultTransport.MaxIdleConnsPerHost = 3

		lvClient = &http.Client{Transport: &pvDefaultTransport}
	}

	pvResponse, pvError := http.DefaultClient.Do(pvRequest)
	if pvError != nil {
		return map[string]interface{}{}, fmt.Errorf("[PBR0003] %v, [%v]", pvError, parURL)
	}

	var pvResponseObj map[string]interface{}
	if pvResponse.StatusCode >= 300 {
		if pvResponse != nil {
			pvResponse.Body.Close()
		}

		pvError = fmt.Errorf("[PBR0004] http error: StatusCode=%v, %v, [%v]", pvResponse.StatusCode, pvResponse.Status, parURL)
		pvResponseObj, _ = convert.NewMapFromHttpResponse(pvResponse)
		return pvResponseObj, pvError
	}

	if pvResponse != nil {
		pvResponseObj, pvError = convert.NewMapFromHttpResponse(pvResponse)

		if pvError != nil {
			pvError = fmt.Errorf("[PBR0005] %v", pvError)
		}
	}

	return pvResponseObj, pvError
}
