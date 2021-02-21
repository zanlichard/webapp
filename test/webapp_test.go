package client

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"gitee.com/cristiane/go-common/json"
)

const (
	baseUrlTestAli = "http://39.108.175.129:52002/api"
	baseUrlDev     = "http://47.106.179.56:52002/api"
	baseUrlLocal   = "http://localhost:8080/api"
)
const (
	appVersionCheck       = "/resources/app/check_version"
)

const (
	apiV1 = "/v1"
	apiV2 = "/v2"
)

var apiVersion = apiV2
var qToken     = token_1000008
var baseUrl    = baseUrlDev + apiVersion

//基本接口
func Testwebapp(t *testing.T) {
	t.Run("APP版本获取", TestAppVersionCheck)
}

const (
	SuccessBusinessCode = 0
)

type HttpCommonRsp struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func TestAppVersionCheck(t *testing.T) {
	r := baseUrl + appVersionCheck
	t.Logf("request url: %s", r)
	data := url.Values{}
	data.Set("client_type", "1")
	rsp, err := http.DefaultClient.PostForm(r, data)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("req url: %v status : %v", r, rsp.Status)
	if rsp.StatusCode != http.StatusOK {
		t.Error("StatusCode != 200")
		return
	}
	body, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("req url: %v body : \n%s", r, body)
	var obj HttpCommonRsp
	err = json.Unmarshal(string(body), &obj)
	if err != nil {
		t.Error(err)
		return
	}
	if obj.Code != SuccessBusinessCode {
		t.Errorf("business code != %v", SuccessBusinessCode)
		return
	}
}






