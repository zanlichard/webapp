package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"webapp/application/appinterface"
	"webapp/frame/subsys"
	"webapp/toolkit"

	"gitee.com/cristiane/go-common/json"
)

const (
	baseUrlDev   = "http://192.168.37.131:51001"
	baseUrlLocal = "http://localhost:51001"
)
const (
	appVersionCheck  = "/api/v1/app/check-version"
	versionCheckPath = "/api/v1/app"
	getbasicCfg      = "/admin/get-basic-cfg"
	getdepCfg        = "/admin/get-dep-cfg"
	getlocalCfg      = "/admin/get-local-acl"
)

type HttpCommonRsp struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	SuccessBusinessCode = 0
	serviceId           = "2160034"
	serviceKey          = "h6F2GvOm1Q1pR5ATYbMjUIUyscLiBs3E"
)

func TestRune(t *testing.T) {
	var i int = 122
	c1 := rune(i)
	fmt.Println("98 convert to", string(c1))
}

func TestGetDepCfg(t *testing.T) {
	request := appinterface.DepCfgGetReq{
		IsServicesAll: true,
	}
	jsonStr, err0 := json.MarshalToString(request)
	if err0 != nil {
		t.Error(err0)
		return
	}
	baseUrl := baseUrlDev + getdepCfg
	client := &http.Client{}
	req, err := http.NewRequest("POST", baseUrl, strings.NewReader(jsonStr))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("content-type", "application/json")
	rsp, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("req url: %+v status : %+v", req, rsp.Status)
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
	t.Logf("req url: %v body : \n%s", req, body)
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

func TestGetBasicCfg(t *testing.T) {
	request := appinterface.BasicCfgGetReq{
		"mongo",
	}
	jsonStr, err0 := json.MarshalToString(request)
	if err0 != nil {
		t.Error(err0)
		return
	}
	baseUrl := baseUrlDev + getbasicCfg
	client := &http.Client{}
	req, err := http.NewRequest("POST", baseUrl, strings.NewReader(jsonStr))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("content-type", "application/json")
	rsp, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("req url: %+v status : %+v", req, rsp.Status)
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
	t.Logf("req url: %v body : \n%s", req, body)
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

func TestGetLocalAclCfg(t *testing.T) {
	baseUrl := baseUrlDev + getlocalCfg
	client := &http.Client{}
	req, err := http.NewRequest("POST", baseUrl, strings.NewReader(""))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("content-type", "application/json")
	rsp, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("req url: %+v status : %+v", req, rsp.Status)
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
	t.Logf("req url: %v body : \n%s", req, body)
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

func TestAppVersionCheck(t *testing.T) {
	head := subsys.SubsysHeader{
		CallServiceId: serviceId,
		GroupNo:       "1",
		Interface:     "check_app_version",
		InvokeId:      "d881c11be7ada28f2d7c602a7c3c20bf",
		MsgType:       "request",
		Remark:        "test",
		Timestamp:     "1608105274",
		Version:       "0.0.1",
	}
	appverCheck := appinterface.AppVersionCheckReq{
		ClientType:     1,
		CurrentVersion: "100001",
	}
	param := appinterface.ParamInfo{
		ApiRequest: appverCheck,
	}
	request := appinterface.AppVerCheckMsg{
		Head:  head,
		Param: param,
	}
	jsonStr, err0 := json.MarshalToString(request)
	if err0 != nil {
		t.Error(err0)
		return
	}

	signStr := toolkit.ApiSign(jsonStr, serviceKey)

	baseUrl := baseUrlDev + appVersionCheck

	client := &http.Client{}
	req, err := http.NewRequest("POST", baseUrl, strings.NewReader(jsonStr))
	if err != nil {
		t.Error(err)
		return
	}

	req.Header.Set("HSB-OPENAPI-CALLERSERVICEID", serviceId)
	req.Header.Set("HSB-OPENAPI-SIGNATURE", signStr)
	req.Header.Set("content-type", "application/json")

	rsp, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("req url: %+v status : %+v", req, rsp.Status)
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
	t.Logf("req url: %v body : \n%s", req, body)
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

func TestAppVersionCheckv2(t *testing.T) {
	appverCheck := appinterface.AppVersionCheckReq{
		ClientType:     1,
		CurrentVersion: "100001",
	}
	param := appinterface.ParamInfo{
		ApiRequest: appverCheck,
	}
	functionName := "check-version"
	_, sessionId, err1 := toolkit.GetUniqId(functionName)
	if err1 != nil {
		t.Errorf("get session id failed for:%+v", err1)
		return
	}
	reqver := "0.0.1"

	reqUrl := baseUrlDev + versionCheckPath
	req, err := subsys.SubsysReqSerialize(reqUrl, serviceId, functionName, sessionId, serviceKey, "request", "test", reqver, param)
	if err != nil {
		t.Errorf("generate req failed for:%+v", err)
		return
	}
	rsp, err2 := subsys.SubsysRequest(req)
	if err2 != nil {
		t.Errorf("request failed for:%+v", err2)
		return
	}
	t.Logf("response:%+v", rsp)

	tmpJson, err3 := json.Marshal(rsp.Rsp.Data)
	if err3 != nil {
		t.Errorf("data format failed:%+v", err3)
		return
	}

	appverCheckRsp := appinterface.AppVersionCheckRsp{}
	if err4 := json.Unmarshal(string(tmpJson), &appverCheckRsp); err4 != nil {
		t.Errorf("data format failed:%+v", err4)
		return
	}
	t.Logf("response:%+v", appverCheckRsp)

}

//基本接口
func Testwebapp(t *testing.T) {
	t.Run("APP版本获取", TestAppVersionCheck)
	t.Run("获取基本配置", TestGetBasicCfg)
	t.Run("APP版本获取v2", TestAppVersionCheckv2)
}
