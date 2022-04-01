package main

import (
	"chatRoom/client/clientProcess"
	"fmt"
	"os"
)

//定义两个变量，一个表示用户id, 一个表示用户密码
var userID int
var userPWD string
var userName string

func main() {

	//接收用户选择
	var key int
	var exitFlag bool = false
	//显示菜单
	for {
		fmt.Println("----------------欢迎登陆多人聊天系统------------")
		fmt.Println("\t\t\t 1 登陆聊天室")
		fmt.Println("\t\t\t 2 注册用户")
		fmt.Println("\t\t\t 3 退出系统")
		fmt.Println("\t\t\t 请选择(1-3):")

		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			fmt.Println("登录聊天室")
			fmt.Println("请输入用户的id")
			fmt.Scanf("%d\n", &userID)
			fmt.Println("请输入用户的密码")
			fmt.Scanf("%s\n", &userPWD)

			err := clientProcess.Login(userID, userPWD)
			if err != nil {
				fmt.Println("登录失败")
			} else {
				//fmt.Println("登录成功")
				//exitFlag = true
			}

		case 2:
			fmt.Println("注册用户")
			fmt.Println("注册用户")
			fmt.Println("请输入用户id:")
			fmt.Scanf("%d\n", &userID)
			fmt.Println("请输入用户密码:")
			fmt.Scanf("%s\n", &userPWD)
			fmt.Println("请输入用户名字(nickname):")
			fmt.Scanf("%s\n", &userName)
			//2. 调用UserProcess，完成注册的请求、
			clientProcess.Register(userID, userPWD)
		case 3:
			fmt.Println("退出系统")
			os.Exit(0)
		default:
			fmt.Println("请输入正确数字")
		}
		if exitFlag {
			break
		}
	}
}
