package clientProcess

import (
	"chatRoom/client/utils"
	"chatRoom/common/message"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

//显示登录成功后的界面..
func ShowMenu() {

	fmt.Println("-------恭喜xxx登录成功---------")
	fmt.Println("-------1. 显示在线用户列表---------")
	fmt.Println("-------2. 发送消息---------")
	fmt.Println("-------3. 信息列表---------")
	fmt.Println("-------4. 退出系统---------")
	fmt.Println("请选择(1-4):")
	var key int
	var content string

	//因为，我们总会使用到SmsProcess实例，因此我们将其定义在swtich外部
	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		fmt.Println("显示在线用户列表-")
		outputOnlineUser()
	case 2:
		fmt.Println("你想对大家说的什么:)")
		fmt.Scanf("%s\n", &content)
		//smsProcess := &SmsProcess{}

		//smsProcess.SendGroupMes(content)
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("你选择退出了系统...")
		err := Logout()
		if err != nil {
			break
		}
		os.Exit(0)
	default:
		fmt.Println("你输入的选项不正确..")
	}

}

func serverProcessMes(conn net.Conn) {
	for true {
		fmt.Println("客户端正在等待读取服务器发送数据")
		mes, err := utils.ReadPkg(conn)
		if err != nil {
			fmt.Println("utils.ReadPkg err: ", err)
			return
		}
		//如果读取到了数据,那么就再处理
		fmt.Println("处理服务端数据...", mes)
		//如果读取到消息，又是下一步处理逻辑
		switch mes.Type {

		case message.NotifyUserStatusMesType:
			//1. 取出.NotifyUserStatusMes
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			//2. 把这个用户的信息，状态保存到客户map[int]User中
			updateUserStatus(&notifyUserStatusMes)
			//处理
		case message.SmsMesType: //有人群发消息
			//outputGroupMes(&mes)
		default:
			fmt.Println("服务器端返回了未知的消息类型")
		}
	}
}
