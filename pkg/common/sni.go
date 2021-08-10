package common

import "strings"

func ToWildcardSni(sni string) string {
	var sb strings.Builder
	for i, c := range sni {
		if c == '.' {
			sb.WriteString("*")
			sb.WriteString(sni[i:])
			break
		}
	}
	return sb.String()
}
