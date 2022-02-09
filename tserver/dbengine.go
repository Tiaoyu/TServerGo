package main

import (
	"TServerGo/log"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"xorm.io/xorm"
)

var (
	dbProxy *DBProxy
)

type DBProxy struct {
	Engine         *xorm.Engine
	DriverName     string
	DataSourceName string
}

func init() {
	dbProxy = &DBProxy{}
}

func (db *DBProxy) Init(driverName, dataSourceName string) {
	db.DriverName = driverName
	db.DataSourceName = dataSourceName

	engine, err := xorm.NewEngine(db.DriverName, db.DataSourceName)
	if err != nil {
		log.Debugf("db proxy create failed! Error:%v", err)
	}
	engine.SetTZLocation(time.UTC)
	db.Engine = engine
	log.Debugf("db proxy driver:%v, host:%v", db.DriverName, db.DataSourceName)
}

func (db *DBProxy) Sync() {
	if err := db.Engine.Sync2(new(User)); err != nil {
		log.Errorf("xorm sync User failed!!!")
	}
	if err := db.Engine.Sync2(new(Race)); err != nil {
		log.Errorf("xorm sync Race failed!!!")
	}
}

func (db *DBProxy) Transaction(fun func(session *xorm.Session) (interface{}, error)) error {
	result, err := db.Engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		result, err := fun(session)
		if err != nil {
			return nil, err
		}
		return result, nil
	})

	if err != nil {
		log.Errorf("DB transaction failed! Result:%v err:%v", result, err)
	}
	return err
}
