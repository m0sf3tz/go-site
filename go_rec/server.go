package main

import (
	"fmt"
	//	"log"
	"net"
	//	"os"
	"time"
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

func timeout(t time.Duration, timeout chan bool) {
	time.Sleep(t)
	timeout <- true
}

func main() {
	// Start the IPC listner, TCP clients will hook into the IPC core
	// to communicate to the backend
	go mq_listner()

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
				go Client_handler(conn)
				fmt.Println("New client attached") //TODO: print client TCP
			}
		}
	}
}
