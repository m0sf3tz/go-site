package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func create_ipc_socket(ipc_socket_name string) net.PacketConn {
	full_socket_path := "nil"

	full_socket_path = SOCKET_PATH + ipc_socket_name

	if err := os.RemoveAll(full_socket_path); err != nil {
		log.Fatal("Could not remove unix-socket")
	}

	l, err := net.ListenPacket("unixgram", full_socket_path) //unixgram==DATAGRAM
	if err != nil {
		fmt.Println("linux socket listen error:", err)
		panic(0)
	}

	return l
}

func ipc_wrangler(cs *Client_state) {
	// wait for the the device to send us it's id,
	// only then can we set up the IPC (the fileis named
	// after the device_id)

	device_id := <-cs.sync
	conn := create_ipc_socket(device_id)
	go ipc_reader(conn, cs)
	//	go ipc_writter(conn, cs)

}

func ipc_reader(conn net.PacketConn, cs *Client_state) {
	ipc_packet := make([]byte, PACKET_LEN_MAX)

	for {
		n, _, err := conn.ReadFrom(ipc_packet)
		if err != nil {
			logger(PRINT_NORMAL, "IPC socket read failed")
		}
		logger(PRINT_DEBUG, "IPC_READER read a packet")
		p := packet_unpack(ipc_packet[0:n])
		cs.ipc_to_tcp_writer_chan <- p
	}
}

func setup_ipc(p Packet, cs *Client_state) {
	if p.Packet_type != HELLO_WORLD_PACKET {
		logger(PRINT_FATAL, "first recieved packet was not a hello-world packet!")
	}
	device_id := "nil"

	// Get the deviceID
	if IPC_TEST_MODE == true {
		device_id = IPC_TEST_SOCKET_NAME
	} else {
		device_id = string(get_packet_device_id(p))
		panic(0)
	}

	cs.sync <- device_id
	time.Sleep(time.Millisecond * 100)
}