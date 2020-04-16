package util

// DigitsCount ...
func DigitsCount(number int64) int64 {
	count := int64(0)
	for number != 0 {
		number /= 10
		count++
	}
	return count

}
