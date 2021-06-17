package dbproxy

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var (
	instance *DBProxy
)

type DBProxy struct {
	Engine         *xorm.Engine
	DriverName     string
	DataSourceName string
}

func init() {
	instance = &DBProxy{}
}
func Instance() *DBProxy {
	return instance
}

func (db *DBProxy) Init(driverName, dataSourceName string) {
	db.DriverName = driverName
	db.DataSourceName = dataSourceName

	engine, err := xorm.NewEngine(db.DriverName, db.DataSourceName)
	if err != nil {
		log.Fatalf("dbproxy create failed! Error:%v", err)
	}
	engine.SetTZLocation(time.UTC)
	db.Engine = engine
	log.Printf("dbproxy driver:%v, host:%v", db.DriverName, db.DataSourceName)
}

func (db *DBProxy) Sync() {
	db.Engine.Sync2(new(User))
	db.Engine.Sync2(new(Race))
}

func (db *DBProxy) Transaction(fun func(session *xorm.Session) (interface{}, error)) {
	db.Engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		result, err := fun(session)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}
