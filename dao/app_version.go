package dao

import (
	"webapp/storage"
	"webapp/model/mysql"
)

func GetAppVersionRecord(clientType int8) (*mysql.AppVersion,error)  {
	var versionRecord mysql.AppVersion
	err := storage.GDb.Table(mysql.TABLEAPPVERSION).Where("client_type=?",clientType).Where("is_delete = 0 AND status = 1").Order("Fid DESC").Limit(1).Find(&versionRecord).Error
    if err != nil {
        return nil,err
	}
	return &versionRecord,nil

}
