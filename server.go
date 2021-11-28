package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	// 在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	// 消息广播
	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 监听Message广播消息channel的goroutine, 一旦有消息就全部发给在线user
func (s *Server) ListenMessager() {
	for {
		msg := <-s.Message
		// 一旦有数据, 把msg 发送给全部在线user
		s.mapLock.Lock()
		for _, u := range s.OnlineMap {
			u.C <- msg
		}
		s.mapLock.Unlock()
	}
}
func (s *Server) BoradCast(u *User, msg string) {
	sendMsg := "[" + u.Addr + "]" + u.Name + ": " + msg
	s.Message <- sendMsg
}

func (s *Server) hanler(conn net.Conn) {
	// 当前链接的任务
	// fmt.Println("链接建立成功!")
	// 用户上线, 将用户加入到onLineMap中
	user := NewUser(conn, s)
	user.Online()

	// 接收客户端上线的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err: ", err)
				return
			}
			user.DoMessage(string(buf[:n-1]))
		}
	}()
	// 当前handler阻塞
	select {}
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

	// 启动监听
	go s.ListenMessager()

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
