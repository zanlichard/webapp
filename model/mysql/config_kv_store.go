package mysql

import "time"

type ConfigKvStore struct {
	Id          uint32    `xorm:" pk 'id'"`
	ConfigKey   string    `xorm:"'config_key' "`   // 配置键
	ConfigValue string    `xorm:"'config_value' "` // 配置值
	Prefix      string    `xorm:"'prefix' "`       // 配置前缀
	Suffix      string    `xorm:"'suffix' "`       // 配置后缀
	Status      uint8     `xorm:"'status' "`       // 是否启用 1是 0否
	IsDelete    uint8     `xorm:"'is_delete' "`    // 是否删除 1是 0否
	CreateTime  time.Time `xorm:"'create_time' "`  // 创建时间
	UpdateTime  time.Time `xorm:"'update_time' "`
}