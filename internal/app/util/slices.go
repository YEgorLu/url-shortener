package util

func Intersects[TType comparable](sliceA []TType, sliceB []TType) bool {

	if len(sliceA) == 0 || len(sliceB) == 0 {
		return false
	}
	values := make(map[TType]bool)
	for _, v := range sliceA {
		values[v] = true
	}
	for _, v := range sliceB {
		if _, ok := values[v]; ok {
			return true
		}
	}
	return true
}
