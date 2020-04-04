package main

import "bytes"
import "encoding/binary"
import "log"

func get_packet_len(packet_type int8) int {
	ret := 0
	switch packet_type {
	case DATA_PACKET:
		ret = DATA_PACKET_SIZE
	case CMD_PACKET:
		ret = CMD_PACKET_SIZE
	case INTERNAL_ACK_PACKET:
		ret = ACK_PACKET_SIZE
	case LOGIN_PACKET:
		ret = LOGIN_PACKET_SIZE
	default:
		log.Fatal("Unknown packet type recieved: ", packet_type)
	}
	return ret
}

// Converts a golang Packet representation
// into a "c" packet
func packetizer(packet Packet) []byte {
	buf := new(bytes.Buffer)

	err1 := binary.Write(buf, binary.LittleEndian, int8(packet.packet_type))
	err2 := binary.Write(buf, binary.LittleEndian, packet.transaction_id)
	err3 := binary.Write(buf, binary.LittleEndian, packet.consumer_ack_req)
	err4 := binary.Write(buf, binary.LittleEndian, packet.data)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		log.Fatal("binary.Write failed - errors are as follows", err1, err2, err3, err4)
	}
	fmt.Println(buf.Bytes())

	return buf.Bytes()
}
