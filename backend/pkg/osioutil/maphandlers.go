package osioutil

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetObjectPathArrayValue(parData map[string]interface{}, parPath []string) (interface{}, error) {

	fnGetObjValue := func(parNode interface{}, parKey string) (pvReturn interface{}, pvError error) {

		var pvOk bool

		switch parNode.(type) {
		case map[string]interface{}:
			if pvReturn, pvOk = parNode.(map[string]interface{})[parKey]; !pvOk {
				pvError = errors.New("key '<0>' is not found")
			}

		case []interface{}:
			if i, pvError := strconv.ParseInt(parKey, 10, 64); pvError == nil {
				pvArray := parNode.([]interface{})
				if int(i) < len(pvArray) {
					pvReturn = pvArray[i]
				} else {
					pvError = errors.New("Index out of bounds for key '<0>'")
				}
			}

		default:
			pvError = errors.New("no key with name '<0>' at this level")
		}

		return
	}

	pvCheckedPath := ""

	var pvObj interface{} = parData
	var pvError error

	for _, pvKey := range parPath {
		if pvCheckedPath != "" {
			pvCheckedPath += "."
		}
		pvCheckedPath += pvKey

		if pvObj, pvError = fnGetObjValue(pvObj, pvKey); pvError != nil {
			pvError = errors.New(strings.Replace(pvError.Error(), "<0>", pvCheckedPath, -1))
			break
		}
	}

	return pvObj, pvError
}

func GetObjectPathValue(parData map[string]interface{}, parPath string) (interface{}, error) {
	pvPathParts := strings.Split(parPath, ".")

	return GetObjectPathArrayValue(parData, pvPathParts)
}

func SetObjectPathValue(parData map[string]interface{}, parPath string, parValue interface{}) error {
	pvPathParts := strings.Split(parPath, ".")

	return SetObjectPathArrayValue(parData, pvPathParts, parValue)
}

func SetObjectPathArrayValue(parData map[string]interface{}, parPath []string, parValue interface{}) error {

	fnGetObjValue := func(parNode interface{}, parKey string) (pvReturn interface{}, pvError error) {

		var pvOk bool

		switch parNode.(type) {
		case map[string]interface{}:
			if pvReturn, pvOk = parNode.(map[string]interface{})[parKey]; !pvOk {
				pvError = errors.New("key '<0>' is not found")
			}

		case []interface{}:
			if i, pvError := strconv.ParseInt(parKey, 10, 64); pvError == nil {
				pvArray := parNode.([]interface{})
				if int(i) < len(pvArray) {
					pvReturn = pvArray[i]
				} else {
					pvError = errors.New("Index out of bounds for key '<0>'")
				}
			}

		default:
			pvError = errors.New("no key with name '<0>' at this level")
		}

		return
	}

	pvCheckedPath := ""

	var pvObj interface{} = parData
	var pvError error

	for i, pvKey := range parPath {
		if pvCheckedPath != "" {
			pvCheckedPath += "."
		}
		pvCheckedPath += pvKey

		if i+1 == len(parPath) {
			pvObj.(map[string]interface{})[pvKey] = parValue
		} else if pvObj, pvError = fnGetObjValue(pvObj, pvKey); pvError != nil {
			pvError = errors.New(strings.Replace(pvError.Error(), "<0>", pvCheckedPath, -1))
			break
		}
	}

	pvObj = parValue

	return pvError
}

func SaveJSONMapToFile(parData map[string]interface{}, parFileName string, parGZIPCompress bool) (pvError error) {
	pvDataBytes, pvError := json.Marshal(parData)
	if pvError != nil {
		return
	}

	if parGZIPCompress {
		pvDataBytes, pvError = GzipCompressData(pvDataBytes)
		if pvError != nil {
			return
		}
	}

	pvError = os.WriteFile(parFileName, pvDataBytes, 0644)
	if pvError != nil {
		return
	}

	return
}

func SaveJSONHTTPStreamToFile(parData io.Reader, parFileName string, parGZIPCompressedData, parGZIPCompress bool, parStampDataKey string) (pvError error) {
	pvDataBytes, pvError := io.ReadAll(parData)
	if pvError != nil {
		return
	}

	if parGZIPCompressedData {
		pvDataBytes, pvError = GzipDecompressData(pvDataBytes)
		if pvError != nil {
			return
		}
	}

	if parStampDataKey != "" {
		var pvDocument map[string]interface{}
		pvError = json.Unmarshal(pvDataBytes, &pvDocument)
		if pvError != nil {
			return
		}

		pvDocument[parStampDataKey] = time.Now()
		return SaveJSONMapToFile(pvDocument, parFileName, parGZIPCompress)
	}

	if parGZIPCompress {
		pvDataBytes, pvError = GzipCompressData(pvDataBytes)
		if pvError != nil {
			return
		}
	}

	pvError = os.WriteFile(parFileName, pvDataBytes, 0644)
	if pvError != nil {
		return
	}

	return
}

func LoadJSONMapFromFile(parFileName string, parGZIPCompress bool) (rtnData map[string]interface{}, rtnError error) {
	b, rtnError := os.ReadFile(parFileName)
	if rtnError != nil {
		return
	}

	if parGZIPCompress {
		b, rtnError = GzipDecompressData(b)
		if rtnError != nil {
			return
		}
	}

	rtnError = json.Unmarshal(b, &rtnData)
	if rtnError != nil {
		return
	}

	return
}

func LoadJSONArrayMapFromFile(parFileName string, parGZIPCompress bool) (rtnData []map[string]interface{}, rtnError error) {
	b, rtnError := os.ReadFile(parFileName)
	if rtnError != nil {
		return
	}

	if parGZIPCompress {
		b, rtnError = GzipDecompressData(b)
		if rtnError != nil {
			return
		}
	}

	rtnError = json.Unmarshal(b, &rtnData)
	if rtnError != nil {
		return
	}

	return
}
