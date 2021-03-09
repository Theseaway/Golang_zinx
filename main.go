// main.go
package main

func main() {
	server := Newserver("127.0.0.1", 9555)
	server.Start()
}
