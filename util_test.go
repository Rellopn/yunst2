package yunst2

import "testing"

func TestSign(t *testing.T) {
	EncodingStr := "test"
	SetPfxPath("./1902271423530473681.pfx")
	SetPfxPwd("123456")
	sign, err := Sign(EncodingStr)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(sign)
	}
}
