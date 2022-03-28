package main

import (
	"chatRoom/server/model"
	"chatRoom/server/severProcess"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"net"
)

//定义一个全局的pool
var pool *redis.Pool

func init() {
	pool = &redis.Pool{
		MaxIdle:     8,   //最大空闲链接数
		MaxActive:   0,   //表示和数据库的最大链接数， 0表示没有限制
		IdleTimeout: 100, //最大空闲时间
		Dial: func() (redis.Conn, error) { //初始化链接的代码，链接哪个ip的redis
			return redis.Dial("tcp", "192.168.137.141:6379")
		},
	}
	model.UDAO = model.NewUserDAO(pool)
}
func main() {
	//创建监听端口
	fmt.Println("服务器监听端口:8888")
	listen, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	defer listen.Close()

	for {
		fmt.Println("等待客户端来链接服务器.....")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}

		//一旦链接成功，则启动一个协程和客户端保持通讯。。
		go process(conn)
	}
}

//处理和客户端的通讯
func process(conn net.Conn) {
	//这里需要延时关闭conn
	defer conn.Close()

	//这里调用总控, 创建一个
	processor := &severProcess.Processor{
		Conn: conn,
	}
	err := processor.ServerProcess()
	if err != nil {
		fmt.Println("客户端和服务器通讯协程错误err: ", err)
		return
	}
}
