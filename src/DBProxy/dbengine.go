package DBProxy

import (
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

type DBEngine struct {
	engine *xorm.Engine
}

var instance *DBEngine

func (e *DBEngine) Init() (err error) {
	e.engine, err = xorm.NewEngine("mysql", "root:#include0@tcp(127.0.0.1:3306)/game?charset=utf8")
	if err != nil {
		return
	}
	e.engine.ShowSQL(true)
	// 生成表结构
	e.engine.Sync2(new(User), new(Role), new(Race))

	return nil
}

func GetInstance() *DBEngine {
	if instance == nil {
		instance = new(DBEngine)
	}
	return instance
}
