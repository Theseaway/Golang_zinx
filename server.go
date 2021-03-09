// server.go
package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

//新建一个端口给服务器服务
func Newserver(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println("Listening on Address:  ", this.Ip, "  Port: ", this.Port)
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
