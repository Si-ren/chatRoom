package utils

import (
	"chatRoom/common/message"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

// Transfer 这里将这些方法关联到结构体中
type Transfer struct {
	//分析它应该有哪些字段
	Conn net.Conn
	Buf  []byte //这时传输时，使用缓冲
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	this.Buf = make([]byte, 2048)
	fmt.Println("读取客户端发送的数据...")
	//conn.Read 在conn没有被关闭的情况下，才会阻塞
	//如果客户端关闭了 conn 则，就不会阻塞
	pkgLen, err := this.Conn.Read(this.Buf)
	if err != nil {
		fmt.Println(err)
		err = errors.New("read pkg header error")
		return
	}
	fmt.Println(pkgLen)
	//把pkgLen 反序列化成 -> message.Message
	// 技术就是一层窗户纸 &mes！！
	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarsha err: ", err)
		return
	}

	dataLen := len(mes.Data) + len(mes.Type)
	fmt.Println(dataLen, mes.Len)
	if dataLen != mes.Len {
		fmt.Println("pkgLen err: ", err)
		return
	}
	return
}

func (this *Transfer) WritePkg(data []byte) (err error) {
	_, err = this.Conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(bytes) fail: ", err)
		return
	}
	return
}
