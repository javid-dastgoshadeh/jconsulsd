package pkg


func IfSliceContainsString(s []string, e string) bool {
	for _, item := range s {
		if item == e {
			return true
		}
	}
	return false
}
