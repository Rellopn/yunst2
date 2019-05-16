# yunst2.0
通联支付 云商通2.0 golang SDK
---
### How to use
```
package config

import (
	"github.com/Rellopn/yunst2"
)

func InitPay() {
	client := yunst2.NewYunClient("http://test/service/soa?", "1", "1",
		"1", "2.0", "./1.pfx", "./t.cer")
	response, mapRes, err := client.Request("MemberService", "createMember", map[string]interface{}{
		"bizUserId":  "20190515105727",
		"memberType": 2,
		"source":     2,
	})
}
```
