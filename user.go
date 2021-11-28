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

func (u *User) DoMessage(msg string) {
	u.server.BoradCast(u, msg)
}

// 监听当前user channel的方法, 一旦有消息, 就直接发送给客户端
func (u *User) listenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}
