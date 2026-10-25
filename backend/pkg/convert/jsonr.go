package convert

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"mosanda-vacation-attendace/pkg/osioutil"
)

func UnmarshalR(parData []byte, parDestination any) error {

	piResult, pvError := isValidJSONFormat(parData, 0)
	if pvError != nil {
		return pvError
	}

	return castValueToGivenDestination(piResult, parDestination)
}

func isValidJSONFormat(parData []byte, parOffset int) (interface{}, error) {

	pvBlockType := ""
	pvDigits := "0123456789"
	i := 0

	fnIndexInc := func() error {
		i++
		if i == len(parData) {
			return osioutil.GetRefErrorTxt("JSONRUN0005", "invalid JSON string EOF")
		}

		return nil
	}

	fnCheckWhiteSpace := func() error {
		if i >= len(parData) {
			return osioutil.GetRefErrorTxt("JSONRUN0030", "invalid JSON format")
		}

		for !strings.Contains(`":,{}[]-0123456789tfn`, string(parData[i])) {
			if parData[i] == 9 || parData[i] == 10 || parData[i] == 13 || parData[i] == 32 {
				pvError := fnIndexInc()
				if pvError != nil {
					return osioutil.GetRefError("JSONRUN0042", pvError)
				}
			} else {
				return osioutil.GetRefErrorTxt("JSONRUN0001", fmt.Sprintf("invalid whitespace character at positon %v", parOffset+i))
			}
		}

		return nil
	}

	fnCheckString := func() (string, error) {
		if parData[i] != 34 {
			return "", osioutil.GetRefErrorTxt("JSONRUN0002", "missing JSON string starting double qoute")
		}

		pvError := fnIndexInc()
		if pvError != nil {
			return "", osioutil.GetRefError("JSONRUN0043", pvError)
		}

		pvReturn := ""

		for {
			if parData[i] < 32 {
				return "", osioutil.GetRefErrorTxt("JSONRUN0003", fmt.Sprintf("invalid JSON string character at position %v", parOffset+i))
			} else if parData[i] == 92 { // \
				pvError := fnIndexInc()
				if pvError != nil {
					return "", osioutil.GetRefError("JSONRUN0044", pvError)
				}

				if parData[i] == 92 || // \
					parData[i] == 47 || // /
					parData[i] == 34 || // "
					parData[i] == 98 || // b
					parData[i] == 102 || // f
					parData[i] == 110 || // n
					parData[i] == 114 || // r
					parData[i] == 116 { // t

					pvReturn += `\` + string(parData[i])
					pvError := fnIndexInc()
					if pvError != nil {
						return "", osioutil.GetRefError("JSONRUN0045", pvError)
					}

				} else if parData[i] == 117 { // u
					pvError := fnIndexInc()
					if pvError != nil {
						return "", osioutil.GetRefError("JSONRUN0046", pvError)
					}

					pvReturn += `\u`

					for j := 0; j < 4; j++ {
						if !strings.Contains(pvDigits+"ABCDEF", strings.ToUpper(string(parData[i]))) {
							return "", osioutil.GetRefErrorTxt("JSONRUN0055", fmt.Sprintf("invalid hex digit at position %v", parOffset+i))
						}

						pvReturn += string(parData[i])

						pvError := fnIndexInc()
						if pvError != nil {
							return "", osioutil.GetRefError("JSONRUN0054", pvError)
						}
					}
				} else {
					return "", osioutil.GetRefErrorTxt("JSONRUN0004", fmt.Sprintf("invalid JSON string's special character at position %v", parOffset+i))
				}
			} else if parData[i] == 34 { // "
				pvError := fnIndexInc()
				if pvError != nil {
					return "", osioutil.GetRefError("JSONRUN0047", pvError)
				}

				break
			} else {
				pvReturn += string(parData[i])
				pvError := fnIndexInc()
				if pvError != nil {
					return "", osioutil.GetRefError("JSONRUN0041", pvError)
				}
			}
		}

		return pvReturn, nil
	}

	fnCheckValidValueTermination := func() error {
		pvError := fnCheckWhiteSpace()
		if pvError != nil {
			return osioutil.GetRefError("JSONRUN0006", pvError)
		}

		if !(parData[i] == 44 || // ,
			parData[i] == 125 || // }
			parData[i] == 93) { // ]
			return osioutil.GetRefErrorTxt("JSONRUN0007", fmt.Sprintf("invalid number termination at position %v", parOffset+i))
		}

		return nil
	}

	fnCheckNumber := func() (interface{}, error) {
		pvValue := ""
		if parData[i] == 45 { // -
			pvValue += "-"

			pvError := fnIndexInc()
			if pvError != nil {
				return nil, osioutil.GetRefError("JSONRUN0048", pvError)
			}
		}

		pvDot := "."
		for {
			if strings.Contains(pvDigits+pvDot, string(parData[i])) {
				pvValue += string(parData[i])

				if parData[i] == 46 {
					pvDot = ""
				}

				pvError := fnIndexInc()
				if pvError != nil {
					return nil, osioutil.GetRefError("JSONRUN0031", pvError)
				}

			} else if strings.Contains(",}] ", string(parData[i])) || parData[i] == 10 || parData[i] == 13 || parData[i] == 9 {
				break
			} else {
				return nil, osioutil.GetRefErrorTxt("JSONRUN0032", fmt.Sprintf("invalid numeric digit at position %v", parOffset+i))
			}
		}

		pvError := fnCheckValidValueTermination()
		if pvError != nil {
			return nil, osioutil.GetRefError("JSONRUN0008", pvError)
		}

		var pvReturn interface{}

		if strings.Contains(pvValue, ".") { // decimal
			pvValueN, pvError := strconv.ParseFloat(pvValue, 64)
			if pvError != nil {
				return nil, osioutil.GetRefError("JSONRUN0056", pvError)
			}

			pvReturn = pvValueN
		} else {
			pvValueN, pvError := strconv.Atoi(pvValue)
			if pvError != nil {
				return nil, osioutil.GetRefError("JSONRUN0057", pvError)
			}

			pvReturn = pvValueN
		}

		return pvReturn, nil
	}

	fnCheckTrue := func() (bool, error) {
		if parData[i] == 116 && // t
			parData[i+1] == 114 && // r
			parData[i+2] == 117 && // u
			parData[i+3] == 101 { // e
			i += 4
		} else {
			return false, osioutil.GetRefErrorTxt("JSONRUN0033", fmt.Sprintf("invalid boolean true at position %v", parOffset+i))
		}

		pvError := fnCheckValidValueTermination()
		if pvError != nil {
			return false, osioutil.GetRefError("JSONRUN0009", pvError)
		}

		return true, nil
	}

	fnCheckFalse := func() (bool, error) {
		if parData[i] == 102 && // f
			parData[i+1] == 97 && // a
			parData[i+2] == 108 && // l
			parData[i+3] == 115 && // s
			parData[i+4] == 101 { // e
			i += 5
		} else {
			return false, osioutil.GetRefErrorTxt("JSONRUN0034", fmt.Sprintf("invalid boolean false at position %v", parOffset+i))
		}

		pvError := fnCheckValidValueTermination()
		if pvError != nil {
			return false, osioutil.GetRefError("JSONRUN0010", pvError)
		}

		return false, nil
	}

	fnCheckNull := func() (interface{}, error) {
		if parData[i] == 110 && // n
			parData[i+1] == 117 && // u
			parData[i+2] == 108 && // l
			parData[i+3] == 108 { // l
			i += 4
		} else {
			return nil, osioutil.GetRefErrorTxt("JSONRUN0035", fmt.Sprintf("invalid null at position %v", parOffset+i))
		}

		pvError := fnCheckValidValueTermination()
		if pvError != nil {
			return nil, osioutil.GetRefError("JSONRUN0011", pvError)
		}

		return nil, nil
	}

	fnCheckSubObject := func() (interface{}, error) {
		var pvBlockStartChar byte = 123 // {
		var pvBlockEndChar byte = 125   // }

		if parData[i] == 91 { // [
			pvBlockStartChar = 91 // [
			pvBlockEndChar = 93   // ]
		}

		pvStartPosition := i
		pvEndPosition := i
		pvBracesBalance := 1

		for pvBracesBalance > 0 {
			pvError := fnIndexInc()
			if pvError != nil {
				return nil, osioutil.GetRefError("JSONRUN0039", pvError)
			}

			pvEndPosition = i

			if parData[i] == pvBlockStartChar {
				pvBracesBalance++
			} else if parData[i] == pvBlockEndChar {
				pvBracesBalance--
			}
		}

		pvError := fnIndexInc()
		if pvError != nil {
			return nil, osioutil.GetRefError("JSONRUN0049", pvError)
		}

		return isValidJSONFormat(parData[pvStartPosition:pvEndPosition+1], parOffset+pvStartPosition)
	}

	fnCheckValue := func() (interface{}, bool, error) {
		pvError := fnCheckWhiteSpace()
		if pvError != nil {
			return nil, false, osioutil.GetRefError("JSONRUN0015", pvError)
		}

		var pvReturn interface{}

		if parData[i] == 34 {
			pvReturn, pvError = fnCheckString()
			if pvError != nil {
				return nil, false, osioutil.GetRefError("JSONRUN0016", pvError)
			}
		} else if strings.Contains("-"+pvDigits, string(parData[i])) {
			pvReturn, pvError = fnCheckNumber()
			if pvError != nil {
				return nil, false, osioutil.GetRefError("JSONRUN0017", pvError)
			}
		} else if parData[i] == 116 { // t
			pvReturn, pvError = fnCheckTrue()
			if pvError != nil {
				return nil, false, osioutil.GetRefError("JSONRUN0018", pvError)
			}
		} else if parData[i] == 102 { // f
			pvReturn, pvError = fnCheckFalse()
			if pvError != nil {
				return nil, false, osioutil.GetRefError("JSONRUN0019", pvError)
			}
		} else if parData[i] == 110 { // n
			pvReturn, pvError = fnCheckNull()
			if pvError != nil {
				return nil, false, osioutil.GetRefError("JSONRUN0020", pvError)
			}
		} else if parData[i] == 123 || parData[i] == 91 { // { OR [
			pvReturn, pvError = fnCheckSubObject()
			if pvError != nil {
				return nil, false, osioutil.GetRefError("JSONRUN0021", pvError)
			}
		} else {
			return nil, false, osioutil.GetRefErrorTxt("JSONRUN0022", fmt.Sprintf("invalid JSON format at position %v", parOffset+i))
		}

		pvError = fnCheckWhiteSpace()
		if pvError != nil {
			return nil, false, osioutil.GetRefError("JSONRUN0023", pvError)
		}

		pfIsCommaTerminated := false

		if parData[i] == 44 { // ,
			pfIsCommaTerminated = true
			pvError := fnIndexInc()
			if pvError != nil {
				return nil, false, osioutil.GetRefError("JSONRUN0050", pvError)
			}
		}

		return pvReturn, pfIsCommaTerminated, nil
	}

	fnCheckObjectElement := func() (string, interface{}, bool, error) {
		pvError := fnCheckWhiteSpace()
		if pvError != nil {
			return "", nil, false, osioutil.GetRefError("JSONRUN0040", pvError)
		}

		pvKey, pvError := fnCheckString() // check field name
		if pvError != nil {
			return "", nil, false, osioutil.GetRefError("JSONRUN0012", pvError)
		}

		pvError = fnCheckWhiteSpace()
		if pvError != nil {
			return "", nil, false, osioutil.GetRefError("JSONRUN0013", pvError)
		}

		if parData[i] == 58 { // :
			pvError := fnIndexInc()
			if pvError != nil {
				return "", nil, false, osioutil.GetRefError("JSONRUN0051", pvError)
			}
		} else {
			return "", nil, false, osioutil.GetRefErrorTxt("JSONRUN0014", fmt.Sprintf("expected colon : at position %v", parOffset+i))
		}

		pvValue, pfIsCommaTerminated, pvError := fnCheckValue()
		if pvError != nil {
			return "", nil, false, osioutil.GetRefError("JSONRUN0058", pvError)
		}

		return pvKey, pvValue, pfIsCommaTerminated, nil
	}

	fnCheckObject := func() (map[string]interface{}, error) {
		pmReturn := make(map[string]interface{})

		for {

			pvKey, pvValue, pfIsCommaTerminated, pvError := fnCheckObjectElement()
			if pvError != nil {
				return nil, osioutil.GetRefError("JSONRUN0024", pvError)
			}

			pmReturn[pvKey] = pvValue

			if parData[i] == 125 { // }
				if pfIsCommaTerminated {
					return nil, osioutil.GetRefErrorTxt("JSONRUN0025", fmt.Sprintf("invalid comma at position %v", parOffset+i))
				}

				break
			}
		}

		return pmReturn, nil
	}

	fnCheckArrayElement := func() (interface{}, bool, error) {
		return fnCheckValue()
	}

	fnCheckArray := func() ([]interface{}, error) {
		paReturn := []interface{}{}

		for {
			pvValue, pfIsCommaTerminated, pvError := fnCheckArrayElement()
			if pvError != nil {
				return nil, osioutil.GetRefError("JSONRUN0037", pvError)
			}

			paReturn = append(paReturn, pvValue)

			if parData[i] == 93 { // ]
				if pfIsCommaTerminated {
					return nil, osioutil.GetRefErrorTxt("JSONRUN0038", fmt.Sprintf("invalid comma at position %v", parOffset+i))
				}

				break
			}
		}

		return paReturn, nil
	}

	pvError := fnCheckWhiteSpace()
	if pvError != nil {
		return nil, osioutil.GetRefError("JSONRUN0026", pvError)
	}

	if parData[i] == 123 { // {
		pvBlockType = "O"
		pvError := fnIndexInc()
		if pvError != nil {
			return nil, osioutil.GetRefError("JSONRUN0052", pvError)
		}
	} else if parData[i] == 91 { // [
		pvBlockType = "A"
		pvError := fnIndexInc()
		if pvError != nil {
			return nil, osioutil.GetRefError("JSONRUN0053", pvError)
		}
	} else {
		return nil, osioutil.GetRefErrorTxt("JSONRUN0027", "invalid JSON format")
	}

	pvError = fnCheckWhiteSpace()
	if pvError != nil {
		return nil, osioutil.GetRefError("JSONRUN0028", pvError)
	}

	if pvBlockType == "O" && parData[i] == 125 { // }
		return make(map[string]interface{}), nil
	} else if pvBlockType == "A" && parData[i] == 93 { // ]
		return []interface{}{}, nil
	}

	var pvReturn interface{}

	if pvBlockType == "O" { // check sub object
		pvReturn, pvError = fnCheckObject()
		if pvError != nil {
			return nil, osioutil.GetRefError("JSONRUN0029", pvError)
		}
	} else if pvBlockType == "A" { // check sub array
		pvReturn, pvError = fnCheckArray()
		if pvError != nil {
			return nil, osioutil.GetRefError("JSONRUN0036", pvError)
		}
	}

	return pvReturn, nil
}

func castValueToGivenDestination(parResult, parDestination any) error {
	pvDestValue := reflect.ValueOf(parDestination)
	if pvDestValue.Kind() != reflect.Pointer || pvDestValue.IsNil() {
		return osioutil.GetRefErrorTxt("JSONRUN0060", "destination type must be a non-null pointer")
	}

	pvDestValueKind := pvDestValue.Elem().Kind()
	switch pvValue := parResult.(type) {
	case map[string]interface{}:
		if pvDestValueKind == reflect.Struct || pvDestValueKind == reflect.Map {
			mapstructure.Decode(pvValue, parDestination)
		} else if pvDestValueKind == reflect.Interface {

		} else {
			return osioutil.GetRefErrorTxt("JSONRUN0062", fmt.Sprintf("cannot unmarshal object into a variable of type %v", pvDestValueKind.String()))
		}

	case []interface{}:
		if pvDestValueKind == reflect.Array ||
			pvDestValueKind == reflect.Slice {

		} else if pvDestValueKind == reflect.Interface {

		} else {
			return osioutil.GetRefErrorTxt("JSONRUN0061", fmt.Sprintf("cannot unmarshal array into a variable of type %v", pvDestValueKind.String()))
		}

	default:
		return osioutil.GetRefErrorTxt("JSONRUN0059", "invalid result's type")
	}

	return nil
}
