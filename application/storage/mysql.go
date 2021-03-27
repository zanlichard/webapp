package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var GDb *gorm.DB
var err error

//
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
		GDb.LogMode(true)
	}
	return nil
}

func ExitDB() {
	if GDb != nil {
		GDb.Close()
	}
}
