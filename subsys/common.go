package subsys

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"webapp/toolkit"

	"gitee.com/cristiane/go-common/json"
)

type SubsysCommonRsp struct {
	Msg     string      `json:"_errStr"`
	Data    interface{} `json:"_data"` //上层应用定义
	RetCode string      `json:"_errCode"`
	Ret     string      `json:"_ret"`
}

type SubsysHeader struct {
	CallServiceId string `json:"_callerServiceId"`
	GroupNo       string `json:"_groupNo"`
	Interface     string `json:"_interface"`
	InvokeId      string `json:"_invokeId"`
	MsgType       string `json:"_msgType"`
	Remark        string `json:"_remark"`
	Timestamp     string `json:"_timestamps"`
	Version       string `json:"_version"`
}

//基本请求体定义
type SubsysReqBody struct {
	Head  SubsysHeader `json:"_head"`
	Param interface{}  `json:"_param"` //上层应用定义
}

//基本响应体定义
type SubsysRspBody struct {
	Head SubsysHeader    `json:"_head"`
	Rsp  SubsysCommonRsp `json:"_data"`
}

func SubsysGetBadHeader() SubsysHeader {
	return SubsysHeader{
		CallServiceId: "unknown",
		GroupNo:       "-1",
		Interface:     "unknown",
		InvokeId:      "unknown",
		MsgType:       "response",
		Remark:        "unknown",
		Timestamp:     toolkit.ConvertToString(toolkit.GetTimeStamp()),
		Version:       "unknown",
	}
}

func SubsysReqSerialize(serviceDomain string, callerServiceId string, funcName string, sessionId string, signKey string, reqType string, reqRemark string, reqVersion string, msgBody interface{}) (*http.Request, error) {
	timestamp := fmt.Sprintf("%d", toolkit.GetTimeStamp())
	head := SubsysHeader{
		CallServiceId: callerServiceId,
		GroupNo:       "1",
		Interface:     funcName,
		InvokeId:      sessionId,
		MsgType:       reqType,
		Remark:        reqRemark,
		Timestamp:     toolkit.ConvertToString(timestamp),
		Version:       reqVersion,
	}
	request := SubsysReqBody{
		Head:  head,
		Param: msgBody,
	}
	jsonStr, err0 := json.MarshalToString(request)
	if err0 != nil {
		return nil, err0
	}
	signStr := toolkit.ApiSign(jsonStr, signKey)
	reqUrl := serviceDomain + "/" + funcName
	req, err := http.NewRequest("POST", reqUrl, strings.NewReader(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("HSB-OPENAPI-CALLERSERVICEID", callerServiceId)
	req.Header.Set("HSB-OPENAPI-SIGNATURE", signStr)
	req.Header.Set("content-type", "application/json")
	return req, nil
}

func SubsysRequest(req *http.Request) (*SubsysRspBody, error) {
	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, errors.New("http reponse statuscode:" + fmt.Sprintf("%d", rsp.StatusCode))
	}
	body, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}
	var obj SubsysRspBody
	err = json.Unmarshal(string(body), &obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil

}
