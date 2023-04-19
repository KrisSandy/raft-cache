package utils

func ConvertToSet(items []string) map[string]struct{} {
	m := make(map[string]struct{})
	for _, i := range items {
		m[i] = struct{}{}
	}
	return m
}
