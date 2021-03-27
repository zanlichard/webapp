package dao

import (
	"webapp/application/dao/model/mysql"
	"webapp/application/storage"
	. "webapp/logger"
)

/*
insert into t_app_version(Fclient_type,Fbuild_code,Fdownload_url,Fforce_update,Fversion_name,Ftitle,Fcontent,Fremark,Fstatus) values('1','test','http://ztiao.club.com',1,"aaa",'xx','xyz','dx','1');
*/
func GetAppVersionRecord(session string, clientType int8) (*mysql.AppVersion, error) {
	var versionRecord mysql.AppVersion
	err := storage.GDb.Table(mysql.TABLEAPPVERSION).Where("Fclient_type=?", clientType).Where("Fis_delete = 0 AND Fstatus = 1").Order("Fid desc").Limit(1).Find(&versionRecord).Error
	if err != nil {
		return nil, err
	}
	InfoFormat("session:%s record:%+v", session, versionRecord)
	return &versionRecord, nil

}
