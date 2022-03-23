package main

import (
	"../common/message"
	"./utils"
	"encoding/json"
	"fmt"
	"net"
	"unsafe"
)

func login(userID int, userPWD string) (err error) {
	user := message.User{
		ID:      userID,
		UserPWD: userPWD,
	}

	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Json Marshal err: ", err)
	}

	var mes message.Message
	mes.Type = message.LoginMesType
	mes.Data = string(data)
	mes.Len = int(unsafe.Sizeof(message.LoginMesType) + unsafe.Sizeof(data) + unsafe.Sizeof(mes.Len))

	data, err = json.Marshal(mes)

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

	mes, err = utils.ReadPkg(conn)

	var loginResMes message.RegisterResMes
	if loginResMes.Code == 200 {
		fmt.Println("登录成功")
	} else if loginResMes.Code == 500 {
		fmt.Println(loginResMes.Error)
	}
	return
}
