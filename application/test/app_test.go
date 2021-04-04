package test

import (
	"gitee.com/cristiane/go-common/json"
	"testing"
	"webapp/application/appinterface"
	"webapp/frame/subsys"
	"webapp/toolkit"
)

const (
	baseUrlDev          = "http://192.168.37.131:51001"
	SuccessBusinessCode = 0
	apiPath             = "/api/v1/app"
	serviceId           = "2160013"
	serviceKey          = "h6F2GvOm1Q1pR5ATYbMjUIUyscLiBs3E"
)

func TestAppVersionCheck(t *testing.T) {
	appverCheck := appinterface.AppVersionCheckReq{
		ClientType:     1,
		CurrentVersion: "100001",
	}

	functionName := "check-version"
	_, sessionId, err1 := toolkit.GetUniqId(functionName)
	if err1 != nil {
		t.Errorf("get session id failed for:%+v", err1)
		return
	}
	reqver := "0.0.1"

	reqUrl := baseUrlDev + apiPath
	req, err := subsys.SubsysReqSerialize(reqUrl, serviceId, functionName, sessionId, serviceKey, "request", "test", reqver, appverCheck)
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

func TestGetImage(t *testing.T) {
	fileKey := "test0001"
	fileMd5 := toolkit.Md5Digest(fileKey)
	fileSize := 33000

	getImageReq := appinterface.GetImageReq{
		FileSize: int32(fileSize),
		FileKey:  fileKey,
		FileMd5:  fileMd5,
	}

	functionName := "get-image"
	_, sessionId, err1 := toolkit.GetUniqId(functionName)
	if err1 != nil {
		t.Errorf("get session id failed for:%+v", err1)
		return
	}
	reqver := "0.0.1"

	reqUrl := baseUrlDev + apiPath
	req, err := subsys.SubsysReqSerialize(reqUrl, serviceId, functionName, sessionId, serviceKey, "request", "test", reqver, getImageReq)
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

	getImageRsp := appinterface.GetImageRsp{}
	if err4 := json.Unmarshal(string(tmpJson), &getImageRsp); err4 != nil {
		t.Errorf("data format failed:%+v", err4)
		return
	}
	t.Logf("response:%+v", getImageRsp)

}

//基本接口
func TestApplication(t *testing.T) {
	t.Run("APP版本获取", TestAppVersionCheck)
	t.Run("获取图片", TestGetImage)
}
