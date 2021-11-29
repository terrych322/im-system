package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(ip string, port int) *Client {
	client := &Client{
		ServerIp:   ip,
		ServerPort: port,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("net.Dial err: ", err)
		return nil
	}
	client.conn = conn

	return client
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址,默认地址是127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器PORT端口,默认地址端口是8888")
}

func main() {
	// 命令行解析
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>>链接服务器成功!")
		return
	}
	fmt.Println(">>>>>>>>链接服务器成功!!!!")
	select {}
}
