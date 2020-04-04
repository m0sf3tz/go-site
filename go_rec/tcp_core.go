package main

import "fmt"
import "net"
import "io"
import "bytes"
import "encoding/binary"
import "log"

type chunker_state struct {
	curr_buff          int
	pckt_size          int
	current_parse_type int
	chunk_buff         []byte
}

func chunker_reset(state chunker_state) {
	state.chunk_buff = nil
	state.pckt_size = 0
	state.current_parse_type = -1
}

func chunker(state chunker_state, rx []byte, lenght int, rx_chan chan Packet) int {
	if len(rx) == 0 || lenght == 0 {
		return -1
	}

	// this function is "stateless", state is passed to it from tcp_socket_read(..)
	curr_buff := state.curr_buff
	pckt_size := state.pckt_size
	current_parse_type := state.current_parse_type
	chunk_buff := state.chunk_buff

	var rx_processed = 0
	for rx_processed != lenght {
		if current_parse_type == -1 {
			current_parse_type = int(rx[PACKET_TYPE_OFFSET]) //the type is always the first byte in any packet
			pckt_size = get_packet_len(int8(current_parse_type))
		}

		var rx_left = lenght - rx_processed
		read_len := 0

		if rx_left <= (pckt_size - curr_buff) {
			read_len = rx_left
		} else {
			read_len = pckt_size - curr_buff
		}

		chunk_buff = append(chunk_buff, rx[rx_processed:rx_processed+read_len]...)
		curr_buff += read_len
		rx_processed += read_len

		if curr_buff == pckt_size {
			//do something
			fmt.Println("Parsed a slice!")

			// Send packet to client_handler
			//TODO
			//res_packet := Packet {
			//rx_chan <- chunk_buff[0:pckt_size]

			//reset internal structures
			chunker_reset(state)
		}
	}
	return 0
}

// Owns low level TCP reads
func Tcp_socket_read(conn net.Conn, err_chan chan bool, rx_chan chan Packet) {
	// Every connection needs a unique chunker
	// create a zero initilzied state, and then get a backing
	// array for the slice
	state := chunker_state{}
	chunker_reset(state)
	state.chunk_buff = make([]byte, PACKET_LEN_MAX)

	chunk := make([]byte, PACKET_LEN_MAX)

	fmt.Println("starting read server")
	for {
		n, err := conn.Read(chunk)

		if err != nil {
			if err != io.EOF {
				fmt.Println("Got the following error on read %s", err)
			} else {
				fmt.Println("Connection closed by client")
			}
			err_chan <- true
			goto err
			return
		}
		fmt.Println("chunking %d bytes", n)
		chunker(state, chunk[0:n], n, rx_chan)
	}
err:
	chunker_reset(state)
	fmt.Println("Closing TCP read socket!")
	conn.Close()
	err_chan <- true
}

// Owns low level TCP writes
func Tcp_socket_write(conn net.Conn, err_chan chan bool, shutdown chan bool, write_chan chan Packet) {
	var packet Packet

	for {
		select {
		case packet = <-write_chan:
			break
		case <-shutdown:
			goto err
		}

		packet_binary := packetizer(packet)

		written := 0
		l := len(packet_binary)
		for written != l {
			n, err := conn.Write(packet_binary[written:])
			if err != nil {
				fmt.Println("got the following error while writting: ", err)
				goto err
			}
			written = written + n
			fmt.Printf("wrote %d bytes\n", n)
		}
	}

err:
	fmt.Println("Closing TCP write socket!")
	conn.Close()
	err_chan <- true
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
