package common

const wildcard = '*'
const wildcardStr = "*"

func ToWildcardSni(sni string) string {
	if len(sni) == 0 {
		return wildcardStr
	}

	if sni[0] == wildcard {
		return sni
	}

	for i, c := range sni {
		if c == '.' {
			return wildcardStr + sni[i:]
		}
	}

	return wildcardStr
}
