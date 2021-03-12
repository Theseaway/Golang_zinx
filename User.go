// User
package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn_ net.Conn, server_ *Server) *User {

	userAddr := conn_.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn_,
		server: server_,
	}
	//启动监听
	go user.ListenMessage()

	return user
}

//监听当前User channel的方法
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

// 上线
func (this *User) Online() {

	//用户上线，将用户加入到OnlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	this.server.Broadcast(this, "Online Now")
}

// 下线

func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	this.server.Broadcast(this, "Offline Now\n")
}

func (this *User) Sendmsg(msg string) {
	//只发送给本地
	this.conn.Write([]byte(msg))
}

// 发送消息的业务
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "Online...\n"
			this.Sendmsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[0:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.Sendmsg("Name has been used\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.Sendmsg("Name changing successfully\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.Sendmsg("Wrong format to send private message\n")
			return
		}

		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.Sendmsg("该用户名不存在\n")
			return
		}

		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.Sendmsg("Nothing to send\n")
			return
		}
		remoteUser.Sendmsg(this.Name + "send message to you" + content)

	} else {
		this.server.Broadcast(this, msg)
	}

}
