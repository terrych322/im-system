package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (s *Server) hanler(net.Conn) {
	// 当前链接的任务
	fmt.Println("链接建立成功!")
}

func (s *Server) Start() {
	// socket
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen er,err:", err)
		return
	}

	// close listen socket
	defer listener.Close()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err, err:", err)
			continue
		}
		// do handler
		go s.hanler(conn)

	}

}
