package main

import (
	"fmt"
	"net"
	//"os"
	"time"
)

const (
	CONN_HOST = "" /* In Go this means listen to ALL the interfaces */
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
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
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				/* retry timer */
				time.Sleep(time.Second * 5)
				l.Close()
				break
			} else {
				if writeSleep(conn) != nil {
					fmt.Println("Got an error - resetting")
					l.Close()
					break
				}
			}
		}
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) error {

	rx_bite := make([]byte, 9)
	rx := make([]byte, 16)
	for {
		n, err := conn.Read(rx_bite)
		/* take care of the case the connection is gracefully closed */
		if err != nil {
			return err
		}
		rx = append(rx, rx_bite[:n]...)
		if len(rx) == 16 {
			fmt.Println("Got a chunk")
			continue
		}
	}
	// Close the connection when you're done with it.
	conn.Close()
	return nil
}

func chunker([]byte rx) nil{

}

func writeSleep(conn net.Conn) error {
	tx_byte := make([]byte, 9)
	for {
		tx_byte[0] = 1
		conn.Write(tx_byte)
		time.Sleep(time.Millisecond * 700)
	}
}
