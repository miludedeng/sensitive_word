package util

func SliceDuplicClear(s []string) (result map[string]int) {
	result = make(map[string]int)
	s_len := len(s)
	for i := 0; i < s_len; i++ {
		result[s[i]] += 1
	}
	return result
}
