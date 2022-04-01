package severProcess

import (
	"chatRoom/common/message"
	"chatRoom/server/model"
	"chatRoom/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	//字段
	Conn net.Conn
	//增加一个字段，表示该Conn是哪个用户
	UserID int
}

func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	//核心代码...
	//1. 先从mes 中取出 mes.Data ，并直接反序列化成LoginMes
	var registerMes message.User
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}
	//1先声明一个 resMes
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	//2在声明一个 RegisterResMes，并完成赋值
	var registerResMes message.RegisterResMes

	//我们需要到redis数据库去完成注册.
	err = model.UDAO.Register(&registerMes)

	if err != nil {
		if err == model.USERNOTEXITS {
			registerResMes.Code = 500
			registerResMes.Error = err.Error()
		} else if err == model.USERPWDERROR {
			registerResMes.Code = 403
			registerResMes.Error = err.Error()
		} else {
			registerResMes.Code = 505
			registerResMes.Error = "服务器内部错误..."
		}
	} else {
		registerResMes.Code = 200
		fmt.Println("注册成功")
	}

	//3将 loginResMes 序列化
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}

	//4. 将data 赋值给 resMes
	resMes.Data = string(data)

	//5. 对resMes 进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}
	//6. 发送data, 我们将其封装到writePkg函数
	//因为使用分层模式(mvc), 我们先创建一个Transfer 实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}

// NotifyOthersOnlineUser 通知所有在线的用户
func (this *UserProcess) NotifyOthersOnlineUser(userID int) {

	//遍历 onlineUsers, 然后一个一个的发送 NotifyUserStatusMes
	for id, up := range userMgr.onlineUsers {
		//过滤到自己
		if id == userID {
			continue
		}
		//开始通知
		up.NotifyMeOnline(userID)
	}
}

//通知
func (this *UserProcess) NotifyMeOnline(userID int) {

	//组装我们的NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserID = userID
	notifyUserStatusMes.Status = message.UserOnline

	//将notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//将序列化后的notifyUserStatusMes赋值给 mes.Data
	mes.Data = string(data)
	mes.Len = len(mes.Data) + len(mes.Type)
	fmt.Println(mes)
	//对mes再次序列化，准备发送.
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//发送,创建我们Transfer实例，发送
	tf := &utils.Transfer{
		Conn: this.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline err=", err)
		return
	}
}

// ServerProcessLogin 处理登录请求
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	//核心代码...
	//1. 先从mes 中取出 mes.Data ，并直接反序列化成LoginMes
	var loginMes message.User
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}
	//1先声明一个 resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType
	//2在声明一个 LoginResMes，并完成赋值
	var loginResMes message.LoginResMes

	//我们需要到redis数据库去完成验证.
	//1.使用model.MyUserDao 到redis去验证
	_, err = model.UDAO.CheckPWD(loginMes.UserID, loginMes.UserPWD)
	if err != nil {
		if err == model.USERNOTEXITS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.USERPWDERROR {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误..."
		}
	} else {
		loginResMes.Code = 200
		this.UserID = loginMes.UserID
		userMgr.AddOnlineUser(this)
		this.NotifyOthersOnlineUser(this.UserID)
		fmt.Println("登录成功")
	}

	//将当前在线用户的id 放入到loginResMes.UsersId
	//遍历 userMgr.onlineUsers
	for id, _ := range userMgr.onlineUsers {
		loginResMes.UsersID = append(loginResMes.UsersID, id)
	}
	//3将 loginResMes 序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}

	//4. 将data 赋值给 resMes
	resMes.Data = string(data)

	//5. 对resMes 进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}
	//6. 发送data, 我们将其封装到writePkg函数
	//因为使用分层模式(mvc), 我们先创建一个Transfer 实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}
