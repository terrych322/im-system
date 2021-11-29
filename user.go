package main

import (
	"fmt"
	"net"
	"strings"
)

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
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		_, ok := s.OnlineMap[newName]
		if ok {
			fmt.Println("当前用户名被使用")
		} else {
			s.mapLock.Lock()
			delete(s.OnlineMap, u.Name)
			s.OnlineMap[newName] = u
			s.mapLock.Unlock()
			u.Name = newName
			u.sendMessage("更改用户名成功, 新用户名是:" + newName + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 获取对方的用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.sendMessage("当前发送信息格式不正确, 请使用\"to|张三|你好啊\"格式")
			return
		}
		remoteUser, ok := s.OnlineMap[remoteName]
		if !ok {
			u.sendMessage("该用户不存在")
		}
		// 根据用户名 得到对方的user对象
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.sendMessage("无消息内容, 请重发!")
			return
		}
		remoteUser.sendMessage(u.Name + "对你说: " + content)
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
