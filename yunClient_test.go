package yunst2

import (
	"fmt"
	"log"
	"testing"
)

func TestYunClient_Post(t *testing.T) {
	client := NewYunClient("http://localhost:6900/service/soa?", "1", "1",
		"1", "2.0", "./1.pfx", "./t.cer")
	response, mapRes, err := client.Request("MemberService", "createMember", map[string]interface{}{
		"bizUserId":  "20190515105727",
		"memberType": 2,
		"source":     2,
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(mapRes)
	log.Println(response)
}
func TestYunClientSMRZ_Post(t *testing.T) {
	client := NewYunClient("http://116.228.64.55:6900/service/soa?", "1902271423530473681", "123456",
		"1902271423530473681", "2.0", "./1902271423530473681.pfx", "./TLCert-test.cer")
	s, e := EncryptionSI("37010319870821451X")
	if e != nil {
		log.Fatal(e)
	}
	fmt.Println(s)
	response, mapRes, err := client.Request("MemberService", "setRealName", map[string]interface{}{
		"bizUserId":    "GR55061321921597440",
		"isAuth":       true,
		"name":         "董然",
		"identityType": 1,
		"identityNo":   s,
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(mapRes)
	log.Println(response)
}

//GR55036727269527552
func TestYunClientJump_Post(t *testing.T) {
	client := NewYunClient("http://116.228.64.55:6900/service/soa?", "1902271423530473681", "123456",
		"1902271423530473681", "2.0", "./1902271423530473681.pfx", "./TLCert-test.cer")
	response, mapRes, err := client.Request("MemberService", "signContract", map[string]interface{}{
		"bizUserId": "GR55061321921597440",
		"jumpUrl":   "back",
		"backUrl":   "before",
		"source":    2,
	}, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(mapRes)
	log.Println(response)
}
