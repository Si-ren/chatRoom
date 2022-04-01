package model

import (
	"chatRoom/common/message"
	"net"
)

var CurUser CurrentUser

//当前用户
type CurrentUser struct {
	Conn net.Conn
	message.User
}
