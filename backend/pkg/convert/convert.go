package convert

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"strings"

	"reflect"

	"mosanda-vacation-attendace/pkg/osioutil"
)

func AnyToBytes(parData any) []byte {

	fnConvertNumeric := func(parData1 any) []byte {

		pvBuffer := new(bytes.Buffer)
		pvError := binary.Write(pvBuffer, binary.LittleEndian, parData1)
		if pvError != nil {
			return []byte{}
		}

		return pvBuffer.Bytes()
	}

	switch pvValue := parData.(type) {
	case byte, uint16, uint32, uint64, int8, int16, int32, int64, float32, float64:
		return fnConvertNumeric(pvValue)

	case *byte, *uint16, *uint32, *uint64, *int8, *int16, *int32, *int64, *float32, *float64:
		return fnConvertNumeric(pvValue)

	case string:
		return []byte(pvValue)

	case *string:
		return []byte(*pvValue)

	default:

	}

	if reflect.ValueOf(parData).Kind() == reflect.Ptr {
		b, pvError := json.Marshal(parData)
		if pvError != nil {
			return []byte{}
		}

		return b
	}

	if reflect.ValueOf(parData).Kind() == reflect.Struct ||
		reflect.ValueOf(parData).Kind() == reflect.Map ||
		reflect.ValueOf(parData).Kind() == reflect.Array ||
		reflect.ValueOf(parData).Kind() == reflect.Slice {
		b, pvError := json.Marshal(parData)
		if pvError != nil {
			return []byte{}
		}

		return b
	}

	return []byte{}
}

func BytesToByte(parData ...byte) (byte, error) {
	if len(parData) < 1 {
		return 0, osioutil.GetRefErrorTxt("BTOB0001", "invalid byte array length")
	}

	return parData[0], nil
}

func BytesToUint16(parData ...byte) (uint16, error) {
	if len(parData) < 2 {
		return 0, osioutil.GetRefErrorTxt("BTOUI160001", "invalid byte array length")
	}

	return binary.LittleEndian.Uint16(parData), nil
}

func BytesToUint32(parData ...byte) (uint32, error) {
	if len(parData) < 4 {
		return 0, osioutil.GetRefErrorTxt("BTOUI320001", "invalid byte array length")
	}

	return binary.LittleEndian.Uint32(parData), nil
}

func BytesToUint64(parData ...byte) (uint64, error) {
	if len(parData) < 8 {
		return 0, osioutil.GetRefErrorTxt("BTOUI640001", "invalid byte array length")
	}

	return binary.LittleEndian.Uint64(parData), nil
}

func BytesToInt16(parData ...byte) (int16, error) {
	if len(parData) < 2 {
		return 0, osioutil.GetRefErrorTxt("BTOI160001", "invalid byte array length")
	}

	var pvReturn int16
	pvError := binary.Read(bytes.NewReader(parData), binary.LittleEndian, &pvReturn)
	if pvError != nil {
		return 0, osioutil.GetRefError("BTOI160002", pvError)
	}

	return pvReturn, nil
}

func BytesToInt32(parData ...byte) (int32, error) {
	if len(parData) < 4 {
		return 0, osioutil.GetRefErrorTxt("BTOI320001", "invalid byte array length")
	}

	var pvReturn int32
	pvError := binary.Read(bytes.NewReader(parData), binary.LittleEndian, &pvReturn)
	if pvError != nil {
		return 0, osioutil.GetRefError("BTOI320002", pvError)
	}

	return pvReturn, nil
}

func BytesToInt64(parData ...byte) (int64, error) {
	if len(parData) < 4 {
		return 0, osioutil.GetRefErrorTxt("BTOI640001", "invalid byte array length")
	}

	var pvReturn int64
	pvError := binary.Read(bytes.NewReader(parData), binary.LittleEndian, &pvReturn)
	if pvError != nil {
		return 0, osioutil.GetRefError("BTOI640002", pvError)
	}

	return pvReturn, nil
}

func BytesToFloat32(parData ...byte) (float32, error) {
	if len(parData) < 4 {
		return 0, osioutil.GetRefErrorTxt("BTOF320001", "invalid byte array length")
	}

	var pvReturn float32
	pvError := binary.Read(bytes.NewReader(parData), binary.LittleEndian, &pvReturn)
	if pvError != nil {
		return 0, osioutil.GetRefError("BTOF320002", pvError)
	}

	return pvReturn, nil
}

func BytesToFloat64(parData ...byte) (float64, error) {
	if len(parData) < 8 {
		return 0, osioutil.GetRefErrorTxt("BTOF640001", "invalid byte array length")
	}

	var pvReturn float64
	pvError := binary.Read(bytes.NewReader(parData), binary.LittleEndian, &pvReturn)
	if pvError != nil {
		return 0, osioutil.GetRefError("BTOF640002", pvError)
	}

	return pvReturn, nil
}

func BytesToString(parData []byte) string {
	return string(parData)
}

func BytesToAny[T any](parData []byte) (T, error) {
	var pvReturn T
	pvError := json.Unmarshal(parData, &pvReturn)
	if pvError != nil {
		return pvReturn, pvError
	}

	return pvReturn, nil
}

func SpecialCharsDecode(parText string) string {
	pvReturn := parText

	fnReplace := func(parChar, parHTMLEntityNumber, parHTMLEntityName string) {
		pvReturn = strings.ReplaceAll(pvReturn, parHTMLEntityNumber, parChar)
		pvReturn = strings.ReplaceAll(pvReturn, parHTMLEntityName, parChar)
	}

	fnReplace(`"`, `&#34;`, `&quot;`)
	fnReplace(`'`, `&#39;`, `&apos;`)
	fnReplace(`&`, `&#38;`, `&amp;`)
	fnReplace(`<`, `&#60;`, `&lt;`)
	fnReplace(`>`, `&#62;`, `&gt;`)
	fnReplace(``, `&#160;`, `&nbsp;`)
	fnReplace(`¢`, `&#162;`, `&cent;`)
	fnReplace(`£`, `&#163;`, `&pound;`)
	fnReplace(`¤`, `&#164;`, `&curren;`)
	fnReplace(`¥`, `&#165;`, `&yen;`)
	fnReplace(`§`, `&#167;`, `&sect;`)
	fnReplace(`©`, `&#169;`, `&copy;`)
	fnReplace(`«`, `&#171;`, `&laquo;`)
	fnReplace(`¬`, `&#172;`, `&not;`)
	fnReplace(`®`, `&#174;`, `&reg;`)
	fnReplace(`°`, `&#176;`, `&deg;`)
	fnReplace(`±`, `&#177;`, `&plusmn;`)
	fnReplace(`»`, `&#187;`, `&raquo;`)
	fnReplace(`¼`, `&#188;`, `&frac14;`)
	fnReplace(`½`, `&#189;`, `&frac12;`)
	fnReplace(`¾`, `&#190;`, `&frac34;`)
	fnReplace(`×`, `&#215;`, `&times;`)
	fnReplace(`÷`, `&#247;`, `&divide;`)

	// if pvReturn != parText {
	// 	fmt.Sprint()
	// }

	return pvReturn
}
