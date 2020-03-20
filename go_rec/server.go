package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var Error_channel chan bool

const SockAddr = "/tmp/echo.sock"

var Tcp_writter_ack chan bool
var Tcp_writter chan Message

var Tcp_core_write_ack chan bool
var Tcp_core_write chan Message

const (
	CONN_HOST = "" /* In Go this means listen to ALL the interfaces */
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

var Device_connected bool = false

var Valid_counter uint16 = 0 // This value is incremented per connection
// and the point of it is to make sure only
// the latest connection can send commands

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

type Message struct {
	Marker int
	data   []byte
}

var writter_q chan Message

func timeout(t time.Duration, timeout chan bool) {
	time.Sleep(t)
	timeout <- true
}

func tcp_core_writter() {

	for {
		x := <-Tcp_writter

		if Device_connected == false {
			fmt.Println("wrote to unconnected device")
			Tcp_writter_ack <- false
			break
		}

		to := make(chan bool, 1)
		go timeout(time.Millisecond*250, to)

		select {
		case Tcp_core_write <- x:
			fmt.Println("wrote to tcp core")
		case <-to:
			fmt.Println("timed out writting to tcp core!")
		}

		fmt.Println(x)

		go timeout(time.Millisecond*250, to)
		select {
		case <-Tcp_core_write_ack:
		// a read from ch has occurred
		case <-to:
			fmt.Println("timed out!")
		}
	}
}

func init_ipc() {
	Tcp_writter = make(chan Message, 1)
	Tcp_writter_ack = make(chan bool, 1)
	Tcp_core_write_ack = make(chan bool, 1)
	Tcp_core_write = make(chan Message, 1)

	//go tcp_core_reader()
	go tcp_core_writter()
}

func main() {
	// Init IPC
	init_ipc()
	p := unix_ipc()
	go ipc_translation_layer(p)

	Error_channel = make(chan bool)
	// Listen for incoming connections.
	for {
		l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
		if err != nil {
			Device_connected = false
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
				Device_connected = false
				fmt.Println("Error accepting: ", err.Error())
				/* retry timer */
				time.Sleep(time.Second * 5)
				break
			} else {
				Device_connected = true
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
	//go Tcp_read(conn, connection_id)
	go Tcp_write(conn)
	<-Error_channel
	conn.Close()
}

func unix_ipc() net.PacketConn {
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

		var m Message
		m.data = make([]byte, 10)
		Tcp_writter <- m
		fmt.Println(b[0:n])
	}
}
