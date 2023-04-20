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

func TestDecryptionSI(t *testing.T) {
	EncodingStr := "6E98DD0AE88D829100795A5B25FAD02F16816A8E5B8B91CC566EE658B056646F1F8E1C37DFB926CAA8804AA1F5C02883D936AD7FF7FE59DE7D9F3248048F6712BF20159142556F3771EFFBB10ED773C5D9CCB46BC2FB52CA7B53760EE666EC93CDE56231E2A9A7B5397E0BF962B181DA918220EA5F4C6908E5FF333445BDCE8B"
	SetPfxPath("./1902271423530473681.pfx")
	SetPfxPwd("123456")
	GetPair()
	de, err := DecryptionSI(EncodingStr)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(de)
	}
}
