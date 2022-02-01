package validation

import "strings"

func IsExist(input string, stringArray []string) bool {
	input = strings.ToLower(input)
	for _, str := range stringArray {
		if strings.ToLower(str) == input {
			return true
		}
	}
	return false
}

func IsExistInt(input int, intArr []int) bool {
	for _, i := range intArr {
		if i == input {
			return true
		}
	}
	return false
}
