package message

// 消息类型分类
const (
	LoginMesType            = "LoginMes"
	LogoutMesType           = "LogoutMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
)

// Message 发送的消息
type Message struct {
	Len  int    `json:"len"`
	Type string `json:"type"`
	Data string `json:"data"`
}

const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

type RegisterMes struct {
	User User `json:"user"` //类型就是User结构体.
}

// RegisterResMes 服务端返回注册信息
type RegisterResMes struct {
	Code  int    `json:"code"`  // 返回状态码 400 表示该用户已经占有 200表示注册成功
	Error string `json:"error"` // 返回错误信息
}

// LoginResMes 返回登录信息，包括已经登录的userID列表
type LoginResMes struct {
	Code    int    `json:"code"` // 返回状态码 500 表示该用户未注册 200表示登录成功
	UsersID []int  // 增加字段，保存用户id的切片
	Error   string `json:"error"` // 返回错误信息
}

// NotifyUserStatusMes 配合服务器端推送用户状态变化的消息
type NotifyUserStatusMes struct {
	UserID int `json:"userID"` //用户id
	Status int `json:"status"` //用户的状态
}
