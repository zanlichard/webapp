package subsys

import (
	"encoding/json"
	"testing"
	e "webapp/apperrors"
	"webapp/appinterface"
	"webapp/toolkit"
)

const (
	versionCheckPath = "/api/v1/app"
	serviceId        = "2160034"
	serviceKey       = "h6F2GvOm1Q1pR5ATYbMjUIUyscLiBs3E"
	baseUrlDev       = "http://192.168.163.128:51001"
)

func TestAppVersionCheck(t *testing.T) {
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
	req, err := SubsysReqSerialize(reqUrl, serviceId, functionName, sessionId, serviceKey, "request", "test", reqver, param)
	if err != nil {
		t.Errorf("generate req failed for:%+v", err)
		return
	}

	t.Logf("request:%+v", req)

	rsp, err2 := SubsysRequest(req)
	if err2 != nil {
		t.Errorf("request failed for:%+v", err2)
		return
	}

	if rsp.Code != int(e.RetCode_SUCCESS) {
		t.Errorf("call service failed for:%d", rsp.Code)
		return
	}

	tmpJson, err3 := json.Marshal(rsp.Data)
	if err3 != nil {
		t.Errorf("data format failed:%+v", err3)
		return
	}

	appverCheckRsp := appinterface.AppVersionCheckRsp{}
	if err4 := json.Unmarshal(tmpJson, &appverCheckRsp); err4 != nil {
		t.Errorf("data format failed:%+v", err4)
		return
	}
	t.Logf("response:%+v", appverCheckRsp)

}
