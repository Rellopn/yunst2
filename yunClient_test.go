package yunst2

import (
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
