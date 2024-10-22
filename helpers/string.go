package helpers

func Contains(slice []string, item string) bool {
	for _, str := range slice {
		if str == item {
			return true
		}
	}
	return false
}
