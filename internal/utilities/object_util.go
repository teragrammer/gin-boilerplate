package utilities

func IsStringValueExistOnArray(arr []string, value *string) bool {
	if value == nil {
		return false
	}

	for _, v := range arr {
		if v == *value {
			return true
		}
	}

	return false
}
