package mysql

import (
	"time"

	_ "github.com/jinzhu/gorm"
)

/*
create table t_app_version(
	   Fid int unsigned auto_increment,
	   Fclient_type   tinyint unsigned NOT NULL DEFAULT 0 comment '客户端类型',
	   Fbuild_code    varchar(256) NOT NULL DEFAULT '' comment 'build值',
       Fdownload_url  varchar(256) NOT NULL DEFAULT '' comment '下载地址',
       Fforce_update  tinyint unsigned NOT NULL DEFAULT 0 comment '是否强制升级1是，0否',
       Fversion_name  varchar(256) NOT NULL DEFAULT '' comment '版本名',
       Ftitle         varchar(256) NOT NULL DEFAULT '' comment '发布标题',
       Fcontent       varchar(256) NOT NULL DEFAULT '' comment '发布内容',
       Fremark        varchar(256) NOT NULL DEFAULT '' comment '说明',
       Fstatus        tinyint NOT NULL DEFAULT 0 comment '0未发布，1已发布，2已经撤销',
       Fis_delete     tinyint unsigned NOT NULL DEFAULT 0 comment '是否删除 1是，0否',
       Fcreate_time   timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP() comment '创建时间',
       Fupdate_time   timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP() comment '修改时间',
       primary key(Fid)
	)engine=innodb default charset=utf8 auto_increment=1;

*/

const TABLEAPPVERSION = "t_app_version"

type AppVersion struct {
	Id          uint32    `gorm:"primary_key"`
	ClientType  uint8     `gorm:"column:Fclient_type"`  // 客户端类型
	BuildCode   string    `gorm:"column:Fbuild_code"`   // build值
	DownloadUrl string    `gorm:"column:Fdownload_url"` // 下载地址
	ForceUpdate uint8     `gorm:"column:Fforce_update"` // 是否强制升级1是，0否
	VersionName string    `gorm:"column:Fversion_name"`
	Title       string    `gorm:"column:Ftitle"`
	Content     string    `gorm:"column:Fcontent"`
	Remark      string    `gorm:"column:Fremark"`
	Status      int8      `gorm:"column:Fstatus"`      // 0未发布，1已发布，2已经撤销
	IsDelete    uint8     `gorm:"column:Fis_delete"`   // 是否删除 1是，0否
	CreateTime  time.Time `gorm:"column:Fcreate_time"` //创建时间
	UpdateTime  time.Time `gorm:"column:Fupdate_time"`
}

/*
type AppVersion struct {
	Id          uint32    `xorm:"'Fid' "`
	ClientType  uint8     `xorm:"'Fclient_type' "`  // 客户端类型
	BuildCode   string    `xorm:"'Fbuild_code' "`   // build值
	DownloadUrl string    `xorm:"'Fdownload_url' "` // 下载地址
	ForceUpdate uint8     `xorm:"'Fforce_update' "` // 是否强制升级1是，0否
	VersionName string    `xorm:"'Fversion_name' "`
	Title       string    `xorm:"'Ftitle' "`
	Content     string    `xorm:"'Fcontent' "`
	Remark      string    `xorm:"'Fremark' "`
	Status      int8      `xorm:"'Fstatus' "`      // 0未发布，1已发布，2已经撤销
	IsDelete    uint8     `xorm:"'Fis_delete' "`   // 是否删除 1是，0否
	CreateTime  time.Time `xorm:"'Fcreate_time' "` //创建时间
	UpdateTime  time.Time `xorm:"'Fupdate_time' "`
}

*/
