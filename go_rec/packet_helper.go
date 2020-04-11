package main

import "bytes"
import "encoding/binary"
import "log"
import "fmt"

/**********************************************************
*       					Helpers for Packets
*********************************************************/

func get_packet_payload_len(packet_type int8) int {
	ret := 0
	switch packet_type {
	case DATA_PACKET:
		ret = LARGE_PAYLOAD_SIZE
		break
	case CMD_PACKET:
		ret = SMALL_PAYLOAD_SIZE
		break
	case INTERNAL_ACK_PACKET:
		ret = REASON_SIZE
		break
	case LOGIN_PACKET:
		ret = MEDIUM_PAYLOAD_SIZE
		break
	case HELLO_WORLD_PACKET:
		ret = SMALL_PAYLOAD_SIZE
		break
	default:
		log.Fatal("Error! Unknown packet type recieved: ", packet_type)
	}
	return ret
}

func get_packet_len(packet_type int8) int {
	ret := 0
	switch packet_type {
	case DATA_PACKET:
		ret = DATA_PACKET_SIZE
		break
	case CMD_PACKET:
		ret = CMD_PACKET_SIZE
		break
	case INTERNAL_ACK_PACKET:
		ret = ACK_PACKET_SIZE
		break
	case LOGIN_PACKET:
		ret = LOGIN_PACKET_SIZE
		break
	case HELLO_WORLD_PACKET:
		ret = HELLO_PACKET_SIZE
		break
	default:
		log.Fatal("ERRO! Unknown packet type recieved: ", packet_type)
	}
	return ret
}

// Converts a golang Packet representation
// into a "c" packet
func packet_pack(packet Packet) []byte {
	fmt.Println("a")
	buf := new(bytes.Buffer)

	err1 := binary.Write(buf, binary.LittleEndian, int8(packet.Packet_type))
	err2 := binary.Write(buf, binary.LittleEndian, packet.Transaction_id)
	err3 := binary.Write(buf, binary.LittleEndian, packet.Consumer_ack_req)
	err4 := binary.Write(buf, binary.LittleEndian, packet.Crc)

	if len(packet.Data) < get_packet_payload_len(packet.Packet_type) {
		logger(PRINT_FATAL, "packet smaller than expected")
	}

	var err5 error
	if packet.Data != nil {
		err5 = binary.Write(buf, binary.LittleEndian, packet.Data)
	}
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		log.Fatal("binary.Write failed - errors are as follows", err1, err2, err3, err4)
	}

	return buf.Bytes()
}

// Converts a C packet into a golang packet
func packet_unpack(packed_packet []byte) Packet {
	fmt.Println("b")
	// first, extract a general_packet (ie, no payload)
	// using binary.read, then later we will extract the data
	var temp_packet Packet_general
	buf := bytes.NewBuffer(packed_packet)

	err := binary.Read(buf, binary.LittleEndian, &temp_packet)
	if err != nil {
		log.Fatal("Failed to unpack packet, error: ", err)
	}

	// Load the Data
	// - first see how large the payload is
	packet_out := Packet{}

	packet_out.Packet_type = temp_packet.Packet_type
	packet_out.Transaction_id = temp_packet.Transaction_id
	packet_out.Consumer_ack_req = temp_packet.Consumer_ack_req
	packet_out.Crc = temp_packet.Crc

	payload_len := get_packet_payload_len(packet_out.Packet_type)
	if payload_len == 0 {
		// This packet does not have a payload
		return packet_out
	}

	if len(packed_packet) < get_packet_payload_len(packet_out.Packet_type) {
		logger(PRINT_FATAL, "len of b[] is smaller than len of packet")
	}

	packet_out.Data = make([]byte, payload_len)
	packet_out.Data = packed_packet[PAYLOAD_OFFSET : PAYLOAD_OFFSET+payload_len]
	return packet_out
}

func get_packet_device_id(p Packet) int32 {
	if p.Packet_type != HELLO_WORLD_PACKET {
		log.Fatal("tried to get the deviceID from a non- hello world packet!")
	}
	device_id := binary.LittleEndian.Uint32(p.Data[0:4])
	return int32(device_id)
}

func create_ack_pack(p Packet, reason uint8) Packet {
	r := Packet{}
	r.Packet_type = SERVER_ACK_PACKET
	r.Transaction_id = p.Transaction_id
	r.Data = make([]byte, SMALL_PAYLOAD_SIZE)
	r.Data[PAYLOAD_OFFSET_DEVICE_ID] = reason

	return r
}

/**********************************************************
*       					Helpers for Ipc_packets
*********************************************************/

// Converts a golang Ipc_packet representation
// into a "c" Ipc_packet
func ipc_packet_pack(i Ipc_packet) []byte {
	buf := new(bytes.Buffer)

	_, err1 := buf.Write(packet_pack(i.P))
	err2 := binary.Write(buf, binary.LittleEndian, i.Id)

	if err1 != nil || err2 != nil {
		logger(PRINT_FATAL, "binary.Write failed - errors are as follows", err1, err2)
	}
	return buf.Bytes()
}

// Converts a "c" Ipc_packet representation
// into a golang Ipc_packet
func ipc_packet_unpack(b []byte) Ipc_packet {
	if b == nil {
		logger(PRINT_FATAL, "nil packet given to ipc-packet_unpack")
	}

	ip := Ipc_packet{}
	ip.P = packet_unpack(b)

	//ugly math... ip.Id starts at (TOTAL_LEN-LEN(INT32))
	id_offset := len(b) - 4 // 4 is len(uint32)
	ip.Id = binary.LittleEndian.Uint32(b[id_offset:])

	return ip
}

/*
func ipc_test_packing_unpacking() {
	ip := Ipc_packet{}
	ip.P.Packet_type = 1
	ip.Id = 0xdeadbeef
	ip.P.Data = make([]byte, SMALL_PAYLOAD_SIZE)
	ip.P.Data[SMALL_PAYLOAD_SIZE-1] = 12
	b := ipc_packet_pack(ip)
	ip_unpack := ipc_packet_unpack(b)

	fmt.Println(b)
	fmt.Println(ip_unpack)
}
*/
