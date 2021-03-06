package clientProcess

import (
	"chatRoom/client/model"
	"chatRoom/client/utils"
	"chatRoom/common/message"
	"encoding/json"
	"fmt"
	"net"
)

func Register(userID int, userPWD string) (err error) {
	//构成用户数据并序列化
	user := message.User{
		UserID:  userID,
		UserPWD: userPWD,
	}
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Json Marshal err: ", err)
	}

	//构成要发送的数据message并序列化
	var mes message.Message
	mes.Len = len(message.RegisterMesType) + len(data)
	mes.Type = message.RegisterMesType
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
	mes, err = utils.ReadPkg(conn)
	var RegisterResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &RegisterResMes)
	if RegisterResMes.Code == 200 {
		fmt.Println("注册成功")
	} else {
		fmt.Println(RegisterResMes.Error)
	}
	return
}

// Login 用于用户登录
func Login(userID int, userPWD string) (err error) {
	//构成用户数据并序列化
	user := message.User{
		UserID:  userID,
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
	mes, err = utils.ReadPkg(conn)
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		fmt.Println("登录成功")
		model.CurUser.Conn = conn
		model.CurUser.UserID = userID
		model.CurUser.UserStatus = message.UserOnline
		for _, v := range loginResMes.UsersID {
			//如果我们要求不显示自己在线,下面我们增加一个代码
			if v == userID {
				continue
			}

			fmt.Println("用户id:\t", v)
			//完成 客户端的 onlineUsers 完成初始化
			user := &message.User{
				UserID:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Print("\n\n")
		//需要一个协程不停读取服务器发送的信息
		//该协程保持和服务器端的通讯.如果服务器有数据推送给客户端
		//则接收并显示在客户端的终端.
		go serverProcessMes(conn)
		for {
			ShowMenu()
		}
	} else {
		fmt.Println(loginResMes.Error)
		return
	}
	return

}

func Logout() (err error) {
	var mes message.Message
	mes.Type = message.LogoutMesType

	var userStatus message.NotifyUserStatusMes
	userStatus.Status = message.UserOffline
	userStatus.UserID = model.CurUser.UserID

	data, err := json.Marshal(userStatus)
	if err != nil {
		fmt.Println("Logout json.Marshal err: ", err)
		return
	}
	mes.Data = string(data)
	mes.Len = len(mes.Data) + len(mes.Type)
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("Logout json.Marshal err: ", err)
		return
	}
	err = utils.WritePkg(model.CurUser.Conn, data)
	return err
}
