package util

import (
	"testing"
)

func Test_MD5(t *testing.T) {
	str := "Test string"
	if MD5(str) != "0fd3dbec9730101bff92acc820befc34" {
		t.Error("MD5 sum incorrect")
	}
}
