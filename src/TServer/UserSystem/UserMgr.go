package UserSystem

import (
	"net"
)

type User struct {
	conn   *net.Conn
	RoleId int32
}

var (
	UserMap       map[string]*User
	UserRoleIdMap map[int32]*User
)

func UserLogin(u *User) {
	UserMap[(*u.conn).RemoteAddr().String()] = u
	UserRoleIdMap[u.RoleId] = u
}

func GetUserByAddr(addr string) *User {
	return UserMap[addr]
}
