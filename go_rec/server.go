package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

var Error_channel chan bool

const (
	CONN_HOST = "" /* In Go this means listen to ALL the interfaces */
	CONN_PORT = "3332"
	CONN_TYPE = "tcp"
)

var Valid_counter uint16 = 0 // This value is incremented per connection
// and the point of it is to make sure only
// the latest connection can send commands

const SockAddr = "/tmp/echo.sock"

func set_time_outs(conn *net.Conn) error {
	err1 := (*conn).SetReadDeadline(time.Now().Add(time.Second * 5))
	err2 := (*conn).SetWriteDeadline(time.Now().Add(time.Second * 5))
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}

// Sets up the channels which talk to the tcp core (internally to this program)
// and also sets up the unix domain sockets which talk externally to this program
func init_ipc() net.PacketConn {
	if err := os.RemoveAll(SockAddr); err != nil {
		fmt.Println("Could not clear unix-socket")
		panic(0)
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
	ipc_translation_layer(l)

	Error_channel = make(chan bool)
	// Listen for incoming connections.
	for {
		l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
		if err != nil {
			fmt.Println("Error listening:", err.Error())
			/* strange... retry in a bit */
			time.Sleep(time.Second * 5)
			continue
		}
		// Close the listener when the application closes.
		defer l.Close()
		fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
		for {
			// Listen for an incoming connection.
			conn, err := l.Accept()
			set_time_outs(&conn)

			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				/* retry timer */
				time.Sleep(time.Second * 5)
				break
			} else {
				Valid_counter++
				go handleRequest(conn, Valid_counter) // Will block here
				fmt.Println("")
			}
		}
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, connection_id uint16) {
	// Start the listener and writter goroutines
	go Tcp_read(conn, connection_id)
	//	go Tcp_write(conn, connection_id)
	<-Error_channel
	conn.Close()
}
