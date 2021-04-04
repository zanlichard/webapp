package test

import (
	"fmt"
	"gitee.com/cristiane/go-common/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"webapp/application/appinterface"
)

const (
	baseUrlDev   = "http://192.168.37.131:51001"
	baseUrlLocal = "http://localhost:51001"
)
const (
	getbasicCfg = "/admin/get-basic-cfg"
	getdepCfg   = "/admin/get-dep-cfg"
	getlocalCfg = "/admin/get-local-acl"
)

type HttpCommonRsp struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	SuccessBusinessCode = 0
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
