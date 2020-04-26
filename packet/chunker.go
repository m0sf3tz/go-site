package main

import "fmt"
import "net"
import "io"

const DATA_PACKET = 0
const CMD_PACKET = 1

const SMALL_PAYLOAD_SIZE = 8
const LARGE_PLAYLOAD_SIZE = 512

const DATA_PACKET_SIZE = (1 + LARGE_PLAYLOAD_SIZE) // extra one for "type"
const CMD_PACKET_SIZE = (1 + SMALL_PAYLOAD_SIZE)   // extra one for "type"
const MES_LEN_MAX = (DATA_PACKET_SIZE)

var curr_buff int
var pckt_size int
var current_parse_type = -1
var chunk_buff []byte

var c chan []byte

func chunker(rx []byte, lenght int) int {
	if len(rx) == 0 || lenght == 0 {
		return -1
	}

	var rx_processed = 0
	for rx_processed != lenght {
		if current_parse_type == -1 {
			current_parse_type = int(rx[rx_processed]) //the type is always the first byte in any packet

			switch current_parse_type {
			case DATA_PACKET:
				pckt_size = DATA_PACKET_SIZE
			case CMD_PACKET:
				pckt_size = CMD_PACKET_SIZE
			default:
				return -1
			}
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
			//c <- chunk_buff

			//reset internal structures
			chunk_buff = nil
			curr_buff = 0
			current_parse_type = -1
		}
	}
	return 0
}

func Tcp_read(conn net.Conn, connection_id uint16) {
	fmt.Println("stating read server")

	for {
		chunk := make([]byte, 100)
		n, err := conn.Read(chunk)

		if err != nil {
			if err != io.EOF {
				fmt.Println("Got the following error on read %s", err)
			} else {
				fmt.Println("Connection closed by client")
			}
			Error_channel <- true
			return
		}
		// Before we process, make sure we are the latest valid connection, else just close
		if connection_id != Valid_counter {
			fmt.Println("Unknown - error? got a read on a stale connection")
			Error_channel <- true
			return
		}
		fmt.Println("chunking %d bytes", n)
		chunker(chunk[0:n], n)
	}
}

func Tcp_write(conn net.Conn) {
	for {
		<-Tcp_core_write
		fmt.Println("wrtting to TCP")

		b := make([]byte, 1+2+1+8)
		create_packet(b)

		n, err := conn.Write(b) //TODO: assuming we wrote in one go - not correct!
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("wrote %d bytes\n", n)

		Tcp_core_write_ack <- true
	}
}
