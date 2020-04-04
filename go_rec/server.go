package main

import (
	"fmt"
	//	"log"
	"net"
	//	"os"
	"time"
)

const (
	CONN_HOST = "" /* In Go this means listen to ALL the interfaces */
	CONN_PORT = "3334"
	CONN_TYPE = "tcp"
)

// This is very important, it will prevent stale connections
// from clogging up the server and hogging server cycles/files
func set_time_outs(conn *net.Conn) error {
	err1 := (*conn).SetReadDeadline(time.Now().Add(time.Minute * 1))
	err2 := (*conn).SetWriteDeadline(time.Now().Add(time.Minute * 1))
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}

type Packet struct {
	packet_type      int8
	transaction_id   int16
	consumer_ack_req int8
	data             []byte
}

func timeout(t time.Duration, timeout chan bool) {
	time.Sleep(t)
	timeout <- true
}

/*
func tcp_core_adaptor() {
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
*/

/*
func init_ipc() {
	Tcp_writter = make(chan Message, 1)
	Tcp_writter_ack = make(chan bool, 1)
	Tcp_core_write_ack = make(chan bool, 1)
	Tcp_core_write = make(chan Message, 1)

	//go tcp_core_reader()
	go tcp_core_writter()
}
*/

func main() {
	// Init IPC
	/*
		init_ipc()
		p := unix_ipc()
		go ipc_translation_layer(p)
	*/

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
			//set_time_outs(&conn)

			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				/* retry timer */
				time.Sleep(time.Second * 5)
				break
			} else {
				go client_handler(conn)
				fmt.Println("New client attached") //TODO: print client TCP
			}
		}
	}
}

// Handles incoming requests.
func client_handler(conn net.Conn) {

	// Set up the timeouts
	set_time_outs(&conn)

	// Must create the error channel we will share with the two TCP
	// reader and writter tasks. If any errors occur during read/write/
	// chunking the sub tasks will notify the client hanlder through
	// this error channel

	// the two TCP writer/reader slaves will
	// write into this channel to let client_handler
	// know something is wrong. The tcp_write_shutdown channel
	// is used to let the tcp_socket_writter goroutine to know
	// it is time to shutdown
	err_chan := make(chan bool, 2)
	tcp_write_shutdown := make(chan bool, 1)

	tcp_socket_writer_chan := make(chan Packet, MAX_OUTSTANDING_TCP_CORE_SEND)
	tcp_socket_reader_chan := make(chan Packet, MAX_OUTSTANDING_TCP_CORE_SEND)

	// Start the listener and writter goroutines
	go Tcp_socket_read(conn, err_chan, tcp_socket_reader_chan)
	go Tcp_socket_write(conn, err_chan, tcp_write_shutdown, tcp_socket_writer_chan)

	select {
	case <-err_chan:
		fmt.Println("a TCP error messge was recieved")
		tcp_write_shutdown <- true
		time.Sleep(time.Second * 5)
		goto tcp_close
	}

tcp_close:
	fmt.Println("closing client_handler")
}

/*
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
}*/
