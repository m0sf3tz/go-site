package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func create_ipc_socket(ipc_socket_name string) net.PacketConn {
	full_socket_path := "nil"

	if err := os.RemoveAll(full_socket_path); err != nil {
		log.Fatal("Could not remove unix-socket")
	}

	full_socket_path = SOCKET_PATH + ipc_socket_name

	l, err := net.ListenPacket("unixgram", full_socket_path) //unixgram==DATAGRAM
	if err != nil {
		fmt.Println("linux socket listen error:", err)
		panic(0)
	}

	return l
}

func ipc_wrangler(cs client_state) {
	// wait for the the device to send us it's id,
	// only then can we set up the IPC (the fileis named
	// after the device_id)

	device_id := <-cs.sync
	conn := create_ipc_socket(device_id)
	go ipc_reader(conn, cs)
	//	go ipc_writter(conn, cs)

}

func ipc_reader(conn net.PacketConn, cs client_state) {
	ipc_packet := make([]byte, PACKET_LEN_MAX)

	for {
		conn.ReadFrom(ipc_packet)
		fmt.Println(ipc_packet)
	}
}
