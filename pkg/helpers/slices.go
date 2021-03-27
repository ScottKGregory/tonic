package helpers

func DedupeStringSlice(s []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func Contains(a []string, s string) bool {
	for _, e := range a {
		if e == s {
			return true
		}
	}

	return false
}
