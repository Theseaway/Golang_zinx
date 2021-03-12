// server.go
package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//在线用户的列表
	OnlineMap map[string]*User //map[key] value
	mapLock   sync.RWMutex
	Message   chan string
}

//新建一个端口给服务器服务
func Newserver(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (this *Server) Broadcast(user *User, msg string) {
	send_msg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- send_msg
}

//监听Message的消息
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg //能够群发的关键在这里
		}
		this.mapLock.Unlock()
	}
}
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("Listening on Address:  ", this.Ip, "  Port: ", this.Port)
	//连接建立，创建一个对应于连接的用户
	user := NewUser(conn, this)

	//用户上线，将用户加入到OnlineMap中
	user.Online()
	isLive := make(chan bool)
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				fmt.Printf(user.Name + "   offline now\n")
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Read err: ", err)
				return
			}
			msg := string(buf[0 : n-1]) // 去除 \n
			//用户操作，上线之后发送消息
			user.DoMessage(msg)
			isLive <- true
		}
	}()

	for {
		select {
		case <-isLive:
		//当前用户活跃
		case <-time.After(time.Second * 10):
			user.Sendmsg("you are offline")
			close(user.C)
			conn.Close()
			return
		}
	}
}

func (this *Server) Start() {
	//第一步，socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// close listen socket
	defer listener.Close()

	//启动监听message的goroutine
	go this.ListenMessager()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listener accept err:", err)
			continue
		}

		//do handler
		go this.Handler(conn)
	}

}
