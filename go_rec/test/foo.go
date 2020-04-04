package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Packet struct {
	transaction_id int16
	data           []byte
	packet_type    int
}

func main() {
	buf := new(bytes.Buffer)
	p := Packet{}
	p.data = make([]byte, 100)
	p.transaction_id = 12
	p.packet_type = 1

	err := binary.Write(buf, binary.LittleEndian, p.packet_type)
	err = binary.Write(buf, binary.LittleEndian, p.transaction_id)
	err = binary.Write(buf, binary.LittleEndian, p.data)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	fmt.Printf("% x", buf.Bytes())
}
