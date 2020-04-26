package main

import "net"
import "time"
import "fmt"

func main() {

	// connect to this socket
	conn, _ := net.Dial("tcp", "127.0.0.1:3332")
	b := make([]byte, 9)

	b[0] = 1
	_, err := conn.Write(b)
	fmt.Println(err)
	time.Sleep(time.Second)
	conn.Close()
}
