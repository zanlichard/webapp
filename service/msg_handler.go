package service

import (
	"context"
	"webapp/appinterface"
	"webapp/dao"
	"webapp/errors"
	. "webapp/logger"
)

//获取app可用版本
func GetAppVersion(ctx context.Context, req *appinterface.AppVersionCheckReq) (*appinterface.AppVersionCheckRsp, errors.RetCode) {
	//取最新的1条记录
	result, err := dao.GetAppVersionRecord(req.ClientType)
	if err != nil {
		Logger.Error("GetAppVersion failed for:%+v", err)
		return nil, errors.RetCode_ERR_DB_SERVER
	}
	rsp := new(appinterface.AppVersionCheckRsp)
	rsp.BuildCode = result.BuildCode
	rsp.Content = result.Content
	rsp.DownloadUrl = result.DownloadUrl
	rsp.ForceUpdate = result.ForceUpdate
	rsp.Remark = result.Remark
	rsp.Title = result.Title
	rsp.VersionName = rsp.Title

	Logger.Error("rsp:%+v ", rsp)
	return rsp, errors.RetCode_SUCCESS
}
