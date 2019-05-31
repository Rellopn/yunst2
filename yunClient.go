package yunst2

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var httpClient *http.Client
var caCert []byte

type YunClient struct {
	serverUrl      string
	sysId          string
	pwd            string
	alias          string
	version        string
	path           string
	tlCertPath     string
	signContactUrl string
}

func (y *YunClient) SetSignContactUrl(signContactUrl string) {
	y.signContactUrl = signContactUrl
}

func NewYunClient(serverUrl string, sysId string, pwd string, alias string, version string, path string, tlCertPath string) *YunClient {
	SetPfxPath(path)
	SetPfxPwd(pwd)
	setTlsClient(tlCertPath)
	GetPair()
	return &YunClient{
		serverUrl:  serverUrl,
		sysId:      sysId,
		pwd:        pwd,
		alias:      alias,
		version:    version,
		path:       path,
		tlCertPath: tlCertPath}
}

func setTlsClient(tlPath string) {
	caCerti, err := ioutil.ReadFile(tlPath)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCerti)
	httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
	caCert = caCerti
}

// 加了一个签名源，为页面跳转、异步响应报文验签，签名源串为：sysid + rps + timestamp
func (y *YunClient) Request(service string, method string, params map[string]interface{}, sourceFrom ...int64) (*http.Response, map[string]string, error) {
	up, sign, err := y.buildPostBody(map[string]interface{}{"service": service, "method": method, "param": params})
	if err != nil {
		return nil, nil, err
	}
	trueUrl := y.serverUrl + url.PathEscape(up) + "sign=" + sign
	if sourceFrom != nil && len(sourceFrom) >= 1 && sourceFrom[0] == 1 {
		return nil, map[string]string{"toUrl": y.signContactUrl + url.PathEscape(up) + "sign=" + sign}, nil
	}
	resp, err := httpClient.Post(trueUrl, "application/x-www-form-urlencoded;charset=utf-8", nil)
	if err != nil {
		return nil, nil, err
	}
	var res map[string]string
	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	err = json.Unmarshal(responseBodyBytes, &res)
	if err != nil {
		return nil, nil, err
	}
	// 默认是同步请求
	//if sourceFrom == nil || len(sourceFrom) == 0 {
	if err := verifySign1(res); err != nil {
		return nil, nil, err
		//}
		//} else { //页面跳转、异步响应报文验签
		//	if err := verifySign2(res); err != nil {
		//		return nil, nil, err
		//	}
	}

	return resp, res, nil
}

func (y *YunClient) buildPostBody(params map[string]interface{}) (string, string, error) {
	pBytes, err := json.Marshal(&params)
	if err != nil {
		return "", "", err
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	sign, err := Sign(y.sysId + string(pBytes) + timestamp)
	if err != nil {
		return "", "", err
	}
	sb := strings.Builder{}
	sb.WriteString("sysid=")
	sb.WriteString(y.sysId)
	escape := url.QueryEscape(sign)
	translate := caseTranslate(escape)
	sb.WriteString("&timestamp=")
	sb.WriteString(timestamp)
	sb.WriteString("&v=")
	sb.WriteString(y.version)
	sb.WriteString("&req=")
	sb.WriteString(string(pBytes))
	sb.WriteString("&")
	return sb.String(), translate, nil
}

func caseTranslate(sign string) string {
	runesign := []rune(sign)
	for i := 0; i < len(runesign); i++ {
		if runesign[i] == 37 {
			if runesign[i+1] >= 65 && runesign[i+1] <= 90 {
				runesign[i+1] = runesign[i+1] + 32
			}
			if runesign[i+2] >= 65 && runesign[i+2] <= 90 {
				runesign[i+2] = runesign[i+2] + 32
			}
			i = i + 2
		}
	}
	return string(runesign)
}

func verifySign1(res map[string]string) error {
	if res["signedValue"] == "" {
		return errors.New("signedValue is null")
	}
	return VerifySign(res["signedValue"], res["sign"])
}

// 页面跳转、异步响应报文验签
func VerifySign2(res map[string]string) error {
	if res["sysid"] == "" {
		return errors.New("sysid is null")
	}
	if res["rps"] == "" {
		return errors.New("rps is null")
	}
	if res["timestamp"] == "" {
		return errors.New("timestamp is null")
	}
	willValidateSource := res["sysid"] + res["rps"] + res["timestamp"]
	return VerifySign(willValidateSource, res["sign"])
}
