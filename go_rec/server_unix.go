package main

import (
	"fmt"
	"log"
	"net"
	"os"
	//"time"
)

type Message struct {
	packet_type int
	data        []byte
}

const SockAddr = "/tmp/echo.sock"

// Sets up the channels which talk to the tcp core (internally to this program)
// and also sets up the unix domain sockets which talk externally to this program
func init_ipc() net.PacketConn {
	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal("Could not remove unix-socket")
	}

	l, err := net.ListenPacket("unixgram", SockAddr) //unixgram==DATAGRAM
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
		var b [1024]byte
		n, _, err := l.ReadFrom(b[:])
		if err != nil {
			panic(err)
		}
		if b[0] == 0 {
			fmt.Println("Got a command")
		} else if b[0] == 1 {
			fmt.Println("Got an ACK")
		} else if b[0] == 2 {
			fmt.Println("got data packet")
		} else {
			//TODO: log
			fmt.Println("unexpected command packet")
		}

		fmt.Println(b[0:n])
	}
}

func main() {
	// Init IPC
	l := init_ipc()
	ipc_translation_layer(l)
}
