package service

import (
	"context"
	e "webapp/application/apperrors"
	"webapp/application/appinterface"
	"webapp/application/dao"
	. "webapp/logger"
)

//获取app可用版本
func GetAppVersion(ctx context.Context, sessionId string, req *appinterface.AppVersionCheckReq) (*appinterface.AppVersionCheckRsp, e.RetCode) {
	//取最新的1条记录
	result, err := dao.GetAppVersionRecord(sessionId, req.ClientType)
	if err != nil {
		ErrorFormat("session:%s GetAppVersion failed for:%+v", sessionId, err)
		return nil, e.RetCode_ERR_DB_SERVER
	}
	rsp := new(appinterface.AppVersionCheckRsp)
	rsp.BuildCode = result.BuildCode
	rsp.Content = result.Content
	rsp.DownloadUrl = result.DownloadUrl
	rsp.ForceUpdate = result.ForceUpdate
	rsp.Remark = result.Remark
	rsp.Title = result.Title
	rsp.VersionName = rsp.Title

	DebugFormat("session:%s rsp:%+v ", sessionId, rsp)
	return rsp, e.RetCode_SUCCESS
}
