TServer
==

## 项目结构

需要先设置GOPATH， 所有依赖库会下载到GOPATH中这样不影响项目本身的文件结构

## xorm

初始化数据库

```go
e.engine, err = xorm.NewEngine("mysql", "root:#include0@tcp(127.0.0.1:3306)/game?charset=utf8")
if err != nil {
return
}
e.engine.ShowSQL(true)
// 生成表结构
e.engine.Sync2(new(User), new(Role), new(Race))
```

事务

```go
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
```

## grpc

[https://grpc.io/docs/protoc-installation/](https://grpc.io/docs/protoc-installation/)

## 功能

### 匹配

角色匹配后加入匹配池子
匹配线程每秒从池子中取若干角色
有足够角色后进行创建房间
创建完房间后开始进入房间线程

### 