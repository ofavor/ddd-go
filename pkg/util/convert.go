package util

func Btoi(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func Itob(i int32) bool {
	return i != 0
}
