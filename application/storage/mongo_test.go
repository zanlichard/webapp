package storage

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
	"webapp/toolkit"
)

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

var (
	tTable         = "image"
	tDatabase      = "mytest"
	tHost          = "192.168.37.131:27017"
	tMongoAccount  = "devops"
	tMongoPassWord = "devops"
)

func TestMongoQuery(t *testing.T) {
	err := InitMgo(tHost, tDatabase, tMongoAccount, tMongoPassWord, 20)
	if err != nil {
		t.Errorf("init mongo failed for:%+v ", err)
		return
	}
	mgoSession := GetMgo()
	if mgoSession == nil {
		t.Error("get mongo session failed")
		return
	}
	defer mgoSession.Close()

	fileKey := "test0001"
	fileMd5 := toolkit.Md5Digest(fileKey)
	fileSize := 33000
	q := &bson.M{
		"file_key":  fileKey,
		"file_md5":  fileMd5,
		"file_size": fileSize,
	}
	fileObj := &Image{}
	err = mgoSession.DB(tDatabase).C(tTable).Find(q).One(fileObj)
	if err != nil {
		t.Errorf("find monogo db failed for:%+v", err)
	}
	t.Logf("query mongodb obj:%+v", fileObj)
}

func TestMongoAdd(t *testing.T) {
	err := InitMgo(tHost, tDatabase, tMongoAccount, tMongoPassWord, 20)
	if err != nil {
		t.Errorf("init mongo failed for:%+v ", err)
		return
	}
	mgoSession := GetMgo()
	if mgoSession == nil {
		t.Error("get mongo session failed")
		return
	}
	defer mgoSession.Close()
	fileObj := &Image{}

	fileSize := int32(33000)
	fileKey := "test0001"
	fileMd5 := toolkit.Md5Digest(fileKey)
	fileURL := "http://github.com/zanlichard/master/blob/docs/SystemDesign.png"
	fileName := "SystemDesign.png"
	fileUid := 1308888

	now := int64(time.Now().UnixNano() / 1000000)
	fileObj.FileKey = fileKey
	fileObj.FileName = fileName
	fileObj.FileMd5 = fileMd5
	fileObj.FileURL = fileURL
	fileObj.FileStatus = 0
	fileObj.FileSize = fileSize
	fileObj.UserID = int32(fileUid)
	fileObj.CreateTime = now
	fileObj.ModifyTime = now
	if err = mgoSession.DB(tDatabase).C(tTable).Insert(fileObj); err != nil {
		t.Errorf("insert into mongodb failed for:%+v", err)
	}
}

func TestMongoUpdate(t *testing.T) {
	err := InitMgo(tHost, tDatabase, tMongoAccount, tMongoPassWord, 20)
	if err != nil {
		t.Errorf("init mongo failed for:%+v ", err)
		return
	}
	mgoSession := GetMgo()
	if mgoSession == nil {
		t.Error("get mongo session failed")
		return
	}
	defer mgoSession.Close()

	fileUid := 1308888
	now := int64(time.Now().UnixNano() / 1000000)
	fi := &Image{}
	fi.UserID = int32(fileUid)
	fi.FileStatus = 1 //更新状态

	selector := &bson.M{"_id": fi.UserID}
	count, err := mgoSession.DB(tDatabase).C(tTable).Find(selector).Count()
	if err != nil {
		t.Errorf("query failed for:%+v", err)
		return
	}
	if count == 0 {
		t.Error("query not match record")
		return
	} else {
		fi.ModifyTime = now
		update := &bson.M{"$set": fi}
		if _, err := mgoSession.DB(tDatabase).C(tTable).Upsert(selector, update); err != nil {
			t.Errorf("update table failed for:%+v", err)
		}
	}
}
