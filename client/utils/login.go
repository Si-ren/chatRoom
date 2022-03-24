package utils

import (
	"chatRoom/common/message"
	"encoding/json"
	"fmt"
	"net"
)

// Login 用于用户登录
func Login(userID int, userPWD string) (err error) {
	//构成用户数据并序列化
	user := message.User{
		ID:      userID,
		UserPWD: userPWD,
	}
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Json Marshal err: ", err)
	}

	//构成要发送的数据message并序列化
	var mes message.Message
	mes.Len = len(message.LoginMesType) + len(data)
	mes.Type = message.LoginMesType
	mes.Data = string(data)
	data, err = json.Marshal(mes)

	//建立连接、发送数据、完成后关闭连接
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail", err)
		return
	}

	//从conn读取数据并转换成message
	mes, err = ReadPkg(conn)
	var loginResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		fmt.Println("登录成功")
	} else if loginResMes.Code == 500 {
		fmt.Println(loginResMes.Error)
	}
	return
}
