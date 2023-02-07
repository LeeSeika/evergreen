package mysql

import (
	"evergreen/settings"
	"fmt"

	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Init() error {
	user := settings.Conf.MySQLConfig.User
	password := settings.Conf.MySQLConfig.Password
	host := settings.Conf.MySQLConfig.Host
	port := settings.Conf.MySQLConfig.Port
	dbname := settings.Conf.MySQLConfig.DBName
	maxOpenConn := settings.Conf.MySQLConfig.MaxOpenConn
	maxIdleConn := settings.Conf.MySQLConfig.MaxIdleConn

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbname)
	var err error
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("sql connect failed", zap.Error(err))
		return err
	}
	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdleConn)
	return nil
}

func Close() {
	err := db.Close()
	if err != nil {
		zap.L().Error("mysql close error:%s", zap.Error(err))
	}
}
