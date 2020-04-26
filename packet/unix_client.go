package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

const SockAddr = "/tmp/echo.sock"

// Sets up the channels which talk to the tcp core (internally to this program)
// and also sets up the unix domain sockets which talk externally to this program
func init_ipc() net.PacketConn {

	l, err := net.Dial("unixgram", SockAddr) //unixgram==DATAGRAM
	if err != nil {
		fmt.Println("listen error:", err)
		panic(0)
	}
	return l
}

// Translates external UNIX-DOMAIN socket flow into interal
// go-channel flows
func ipc_translation_layer(l net.PacketConn) {
	for {
		var buf [1024]byte
		n, _, err := l.ReadFrom(buf[:])
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", string(buf[:n]))
	}
}

func main() {
	// Init IPC
	l := init_ipc()
	b := make([]byte, 10)
	b[0] = 1
	l.Write(b, 1)
}
