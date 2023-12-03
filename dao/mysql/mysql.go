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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", user, password, host, port, dbname)
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
		zap.L().Error("mysql close error", zap.Error(err))
	}
}

func BeginTransaction() (*sqlx.Tx, error, func(error)) {
	zap.L().Info("Begin a transaction now...")
	tx, err := db.Beginx()
	if err != nil {
		zap.L().Error("transaction begin failed", zap.Error(err))
		return nil, err, nil
	}
	f := func(err error) {
		if p := recover(); p != nil {
			tx.Rollback()
			zap.L().Error("transaction rollback cause of panic")
			panic(p)
		} else if err != nil {
			// 下面logic业务产生error都会在这里打印
			tx.Rollback()
			zap.L().Error("transaction rollback case of error", zap.Error(err))
		} else {
			tx.Commit()
			zap.L().Info("transaction committed")
		}
	}

	return tx, err, f
}
