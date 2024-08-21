package util

import "testing"

// func Test_AesEncrypt(t *testing.T) {
// 	key := "M19kAJF8B50AdKNp"
// 	str := "t18612345678"
// 	es, _ := AesEncrypt([]byte(key), str)
// 	t.Error(es)
// }

func Test_AesDecrypt(t *testing.T) {
	key := "M19kAJF8B50AdKNp"
	es := "DMZfqlcjvz32udKrgN_MPEeL1oH7-3uYFXR0"
	str, _ := AesDecrypt([]byte(key), es)
	if str != "18612345678" {
		t.Error("Decrypt failed")
	}
}
