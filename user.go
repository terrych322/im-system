package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// 创建一个用户
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	go user.listenMessage()
	return user
}

func (u *User) Online() {
	s := u.server
	s.mapLock.Lock()
	s.OnlineMap[u.Name] = u
	s.mapLock.Unlock()
	s.BoradCast(u, "已上线")
}

func (u *User) Offline() {
	s := u.server
	s.mapLock.Lock()
	delete(s.OnlineMap, u.Name)
	s.mapLock.Unlock()
	s.BoradCast(u, "已下线")
}

// 给当前用户的客户端发送广播
func (u *User) sendMessage(msg string) {
	u.conn.Write([]byte(msg))
}

func (u *User) DoMessage(msg string) {
	s := u.server
	if msg == "who" {
		//查询当前用户都有哪些
		s.mapLock.Lock()
		for _, user := range s.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + "在线...\n"
			u.sendMessage(onlineMsg)
		}
		s.mapLock.Unlock()
	} else {
		s.BoradCast(u, msg)
	}

}

// 监听当前user channel的方法, 一旦有消息, 就直接发送给客户端
func (u *User) listenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}
