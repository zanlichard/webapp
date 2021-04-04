package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
	"webapp/logger"
)

var GDb *gorm.DB
var err error

type SqlLogger struct {
}

func (log *SqlLogger) Print(values ...interface{}) {
	level := values[0]
	source := values[1]
	if level == "sql" {
		sql := values[3].(string)
		logger.InfoFormat("%s--sql:%s", source, sql)
	} else {
		logger.InfoFormat("%+v", values)
	}
}

func InitDB(serverAddr string, user string, pwd string, database string, maxOpen int, maxIdle int, idleTime int, debug bool) error {
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pwd, serverAddr, database)
	GDb, err = gorm.Open("mysql", connStr)
	if err != nil {
		return err
	}
	GDb.DB().SetMaxOpenConns(maxOpen)
	GDb.DB().SetMaxIdleConns(maxIdle)
	GDb.DB().SetConnMaxLifetime(time.Duration(idleTime) * time.Second)

	if debug {
		sqlLog := &SqlLogger{}
		GDb.LogMode(true)
		GDb.SetLogger(sqlLog)
	}
	return nil
}

func ExitDB() {
	if GDb != nil {
		GDb.Close()
	}
}
