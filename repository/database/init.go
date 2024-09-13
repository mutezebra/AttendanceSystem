package database

import (
	"database/sql"
	"fmt"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/config"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var _db *sql.DB

func InitMysql() {
	conf := config.Conf.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=true&loc=Local", conf.UserName, conf.Password, conf.Address, conf.Database, conf.Charset)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.LogrusObj.Panic(err)
	}
	if err = db.Ping(); err != nil {
		log.LogrusObj.Panic(err)
	}
	db.SetMaxOpenConns(150)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(5 * time.Minute)
	_db = db
}
