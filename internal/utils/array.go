package utils

func InStringArray(arr []string, target string) bool {
	if len(arr) == 0 {
		return false
	}
	for _, a := range arr {
		if target == a {
			return true
		}
	}
	return false
}
