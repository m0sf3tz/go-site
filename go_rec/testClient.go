package main

import "net"
import "time"
import "log"
import "fmt"

func main() {
	// connect to TCP server
	fmt.Println("starting TEST client!")

	c, err := net.Dial("tcp", ":"+CONN_PORT)
	if err != nil {
		log.Fatal(err)
	}

	p := Packet{}
	p.Transaction_id = 12
	p.Packet_type = 5
	p.Data = make([]byte, SMALL_PAYLOAD_SIZE)

	packed := packet_pack(p)
	fmt.Println(packed)

	_, err = c.Write(packed)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 100)
}
