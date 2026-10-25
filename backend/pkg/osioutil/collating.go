package osioutil

import "strings"

//GetUnifiedString uinified the given string to a simple searchable form
func GetUnifiedString(parText string) string {
	pvReturn := ""

	paRunes := []rune(parText)

	for _, pvCharCode := range paRunes {

		if pvCharCode >= 65280 && pvCharCode <= 65374 {
			pvCharCode = pvCharCode - 65248
		}

		if pvCharCode >= 1536 && pvCharCode <= 1791 {
			if (pvCharCode >= 1569 && pvCharCode <= 1594) || (pvCharCode >= 1601 && pvCharCode <= 1610) {

				switch pvCharCode {
				case 1570, 1571, 1573: //'آ', 'أ', 'إ'
					pvCharCode = 1575 //'ا'

				case 1572: //'ؤ'
					pvCharCode = 1608 //'و'

				case 1574: //'ئ'
					pvCharCode = 1569 //'ء'

				case 1609: //'ى'
					pvCharCode = 1610 //'ي'

				case 1577: //'ة'
					pvCharCode = 1607 //'ه'
				}

				pvReturn += string(pvCharCode)

			}

			if pvCharCode >= 1632 && pvCharCode <= 1641 {
				pvReturn += string(pvCharCode - 1584)
			}

		} else {
			pvReturn += strings.ToLower(string(pvCharCode))
		}
	}

	return pvReturn
}
