package storage

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"time"
	"webapp/logger"
)

var mongo *mgo.Session
var mgoDataBase string

type MgoLogger struct {
}

func (log *MgoLogger) Output(calldepth int, s string) error {
	logger.InfoFormat("cmd:%s", s)
	return nil
}

// MgoSetMaxpool 设置最大池子
func MgoSetMaxpool(maxpool int) error {
	mongo.SetPoolLimit(maxpool)
	return nil
}

func GetMongoDatabaseName() string {
	return mgoDataBase
}

// InitMgo 初始化mongo
func InitMgo(addr, database, account, passwd string, maxpoolsize int, isDebug bool) (err error) {
	mgoconf := fmt.Sprintf("mongodb://%s?maxPoolSize=%d", addr, maxpoolsize)
	if account != "" {
		mgoconf = fmt.Sprintf("mongodb://%s:%s@%s/%s?maxPoolSize=%d", account, passwd, addr, database, maxpoolsize)
	}

	mgoDataBase = database
	mongo, err = mgo.Dial(mgoconf)
	if err == nil {
		mongo.SetMode(mgo.PrimaryPreferred, true)
	} else {
		return err
	}
	if isDebug {
		mgoLog := &MgoLogger{}
		mgo.SetDebug(true)
		mgo.SetLogger(mgoLog)
	}
	return nil
}

// InitMgoEx 初始化mongo
func InitMgoEx(addr []string, database, user, passwd string, maxpoolsize, timeout int) (err error) {
	di := &mgo.DialInfo{}
	di.Addrs = append(di.Addrs, addr...)
	di.Database = database
	di.Username = user
	di.Password = passwd
	di.PoolLimit = maxpoolsize
	di.FailFast = true
	di.Timeout = time.Duration(timeout) * time.Millisecond
	di.Direct = true
	mongo, err = mgo.DialWithInfo(di)
	if err != nil {
		return err
	}
	mongo.SetMode(mgo.PrimaryPreferred, true)
	mongo.SetSyncTimeout(15 * time.Second)
	mongo.SetSocketTimeout(15 * time.Second)
	go initMgoPool(maxpoolsize)
	return nil
}
func initMgoPool(maxpoolsize int) {
	var sessions []*mgo.Session
	for i := 0; i < maxpoolsize/2; i++ {
		session := mongo.Copy()
		session.Ping()
		sessions = append(sessions, session)
	}

	for k := range sessions {
		sessions[k].Close()
	}
}

// GetMgo 获取mongo会话进行操作
// 注意: 使用完成要close
func GetMgo() *mgo.Session {
	return mongo.Copy()
}

// CloseMgo 关闭mongo
func CloseMgo(mgosesson *mgo.Session) {
	if mgosesson != nil {
		mgosesson.Close()
		return
	}
	mongo.Close()
}
