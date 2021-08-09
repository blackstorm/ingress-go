package common

const EMPTY_STRING = ""

func CheckOrDefault(str *string, val string) string {
	if str == nil {
		return val
	}
	return *str
}
