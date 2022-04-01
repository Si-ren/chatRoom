package severProcess

import (
	"chatRoom/common/message"
	"chatRoom/server/utils"
	"encoding/json"
	"fmt"
)

func (this *UserProcess) ServerProcessLogout(mes *message.Message) (err error) {
	var notifyUserStatusMes message.NotifyUserStatusMes
	err = json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
	if err != nil {
		fmt.Println("ServerProcessLogout json.Unmarshal err: ", err)
		return
	}

	userMgr.DelOnlineUser(notifyUserStatusMes.UserID)
	this.MyNotifyOthersOnlineUser(&notifyUserStatusMes)
	return
}

// MyNotifyOthersOnlineUser  通知所有在线的用户
func (this *UserProcess) MyNotifyOthersOnlineUser(notifyUserStatusMes *message.NotifyUserStatusMes) {

	//遍历 onlineUsers, 然后一个一个的发送 NotifyUserStatusMes
	for _, up := range userMgr.onlineUsers {
		up.MyNotifyMeOnline(notifyUserStatusMes)
	}
}

// MyNotifyMeOnline 通知
func (this *UserProcess) MyNotifyMeOnline(notifyUserStatusMes *message.NotifyUserStatusMes) {

	//组装我们的NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

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
