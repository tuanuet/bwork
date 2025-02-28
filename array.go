package bwork

// DifferentArrays find elements only exist in arr1 (comparable by common field in each arr)
func DifferentArrays[T, M any, K comparable](arr1 []T, arr2 []M, funcGetCompareKey1 func(T) K, funcGetCompareKey2 func(M) K) []T {
	appearedMap := make(map[K]bool)

	onlyInArr1 := make([]T, 0)

	for _, item := range arr2 {
		appearedMap[funcGetCompareKey2(item)] = true
	}

	for _, item := range arr1 {
		_, f := appearedMap[funcGetCompareKey1(item)]
		if !f {
			onlyInArr1 = append(onlyInArr1, item)
		}
	}

	return onlyInArr1
}
