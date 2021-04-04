package dao

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"webapp/application/dao/model/mongodb"
	"webapp/application/storage"
	"webapp/logger"
)

func GetImage(session string, fileSize int32, fileMd5 string, fileKey string) (*mongodb.Image, error) {
	mgoSession := storage.GetMgo()
	if mgoSession == nil {
		return nil, errors.New("get mongo session failed")
	}
	defer mgoSession.Close()
	q := &bson.M{
		"file_key":  fileKey,
		"file_md5":  fileMd5,
		"file_size": fileSize,
	}
	fileObj := &mongodb.Image{}
	dataBase := storage.GetMongoDatabaseName()
	err := mgoSession.DB(dataBase).C(mongodb.TABLEIMAGE).Find(q).One(fileObj)
	if err != nil {
		return nil, err
	}
	logger.InfoFormat("query mongodb:%s obj:%+v", dataBase, fileObj)
	return fileObj, nil

}
