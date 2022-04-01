package utils

import (
	"chatRoom/common/message"
	_ "encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

// ReadPkg 从conn读取数据并判断是否有丢包/粘包,转换成message返回
func ReadPkg(conn net.Conn) (mes message.Message, err error) {

	buf := make([]byte, 8096)
	fmt.Println("读取服户端发送的数据...")
	//conn.Read 在conn没有被关闭的情况下，才会阻塞
	//如果客户端关闭了 conn 则，就不会阻塞

	//读取消息内容
	pkgLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Read err: ", err)
		return
	}

	//把pkgLen 反序列化成 -> message.Message
	err = json.Unmarshal(buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarsha err: ", err)
		return
	}
	fmt.Println(mes)
	pkgLen = len(mes.Data) + len(mes.Type)
	if mes.Len != pkgLen {
		err = errors.New("read pkg body error")
		return
	}

	return
}

func WritePkg(conn net.Conn, data []byte) (err error) {

	//先发送一个长度给对方
	var pkgLen int
	pkgLen = len(data)

	//发送data本身
	n, err := conn.Write(data)
	fmt.Println(n, pkgLen)
	if n != pkgLen || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	return
}
