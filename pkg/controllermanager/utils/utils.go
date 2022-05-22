package utils

func RemoveDuplication(arr []string) []string {
	set := make(map[string]struct{}, len(arr))
	j := 0
	for _, v := range arr {
		_, ok := set[v]
		if ok {
			continue
		}
		set[v] = struct{}{}
		arr[j] = v
		j++
	}

	return arr[:j]
}

func RemoveMultiIndex(orig []string, toRemove []int) []string {
	result := make([]string, 0)
	for idx, str := range orig {
		if len(toRemove) > 0 && toRemove[0] == idx {
			toRemove = toRemove[1:]
		} else {
			result = append(result, str)
		}
	}
	return result
}
