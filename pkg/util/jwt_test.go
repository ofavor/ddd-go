package util

import (
	"testing"
	"time"
)

func Test_JwtEncodeAndDecode(t *testing.T) {
	JwtKey = "i88U2kmkhwq29dkDD2ybb"
	s, err := JwtEncode(
		10,
		"test",
		1,
		time.Hour*10000,
	)
	if err != nil {
		t.Error("JwtEncode failed", err)
	}
	t.Log(s)

	n, err := JwtDecode(s)
	if err != nil {
		t.Error("JwtDecode failed", err)
	}
	if n.UserId != 10 || n.Username != "test" || n.Type != 1 {
		t.Error("JwtDecode failed", n)
	}
	t.Log(n)
}

func Test_JwtEncodeTwiceDifferent(t *testing.T) {
	s1, _ := JwtEncode(
		10,
		"test",
		1,
		time.Second*60,
	)
	s2, _ := JwtEncode(
		10,
		"test",
		1,
		time.Second*60,
	)
	if s1 == s2 {
		t.Error("JwtEncode twice should genearte different token")
	}
}
