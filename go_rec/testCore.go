package main

import "net"
import "time"
import "log"
import "fmt"

func main() {
	// connect to test IPC socket

	fmt.Println("hia")
	c, err := net.Dial("unixgram", SOCKET_PATH+IPC_TEST_SOCKET_NAME)
	if err != nil {
		log.Fatal(err)
	}

	p := Packet{}
	p.Transaction_id = 12
	p.Packet_type = 5
	p.Consumer_ack_req = 1
	p.Data = make([]byte, SMALL_PAYLOAD_SIZE)

	packed := packet_pack(p)
	c.Write(packed)

	time.Sleep(time.Second * 100)
}
