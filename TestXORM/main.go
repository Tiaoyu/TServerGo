package main

import (
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var engine *xorm.Engine

func main() {
	var err error
	engine, err = xorm.NewEngine("mysql", "root:#include0@tcp(127.0.0.1:3306)/game?charset=utf8")
	if err != nil {
		return
	}
	engine.ShowSQL(true)

	// 生成表结构
	engine.Sync2(new(User), new(Role), new(Race))

	// insert例子
	user := new(User)
	user.UserName = "tiaoyu"
	_, err = engine.Insert(user)
	if err != nil {
	}

	if err = MyTransactionOps(); err != nil {
	}
}

// MyTransactionOps 事务例子
func MyTransactionOps() error {
	session := engine.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	user := User{UserName: "tiaoyu1"}
	if _, err := session.Where("user_name=?", "tiaoyu").Update(&user); err != nil {
		return err
	}
	return session.Commit()
}
