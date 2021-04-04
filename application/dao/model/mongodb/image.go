package mongodb

const TABLEIMAGE = "image"

type Image struct {
	UserID     int32  `bson:"_id"` // 唯一
	FileKey    string `bson:"file_key"`
	FileName   string `bson:"file_name"`
	FileMd5    string `bson:"file_md5"`
	FileSize   int32  `bson:"file_size"`
	FileURL    string `bson:"file_url"`
	AppID      int32  `bson:"app_id"`
	FileStatus int32  `bson:"file_status"` // 备用
	CreateTime int64  `bson:"create_time"` // 创建时间
	ModifyTime int64  `bson:"modify_time"` // 修改时间
}
