package main

import "fmt"
import "net"
import "io"

//import "time"

type chunker_state struct {
	packet_parsed_so_far int
	pckt_size            int
	current_parse_type   int
	packet_wip           []byte
}

func chunker_reset(state *chunker_state) {
	state.packet_parsed_so_far = 0
	state.packet_wip = nil
	state.pckt_size = 0
	state.current_parse_type = -1
}

func chunker(state *chunker_state, rx []byte, lenght int, rx_chan chan Packet) int {
	if len(rx) == 0 || lenght == 0 {
		return -1
	}

	var rx_processed = 0
	for rx_processed != lenght {

		// this function is "stateless", state is passed to it from tcp_socket_read(..)
		if state.current_parse_type == -1 {
			// the type is always the first byte in any packet
			// - also, it must be offset to the boundary of the next
			// incomming packet
			state.current_parse_type = int(rx[rx_processed])

			state.pckt_size = get_packet_len(int8(state.current_parse_type))
		}

		var rx_left = lenght - rx_processed
		read_len := 0

		if rx_left <= (state.pckt_size - state.packet_parsed_so_far) {
			read_len = rx_left
		} else {
			read_len = state.pckt_size - state.packet_parsed_so_far
		}

		state.packet_wip = append(state.packet_wip, rx[rx_processed:rx_processed+read_len]...)
		state.packet_parsed_so_far += read_len
		rx_processed += read_len

		if state.packet_parsed_so_far == state.pckt_size {
			fmt.Println("Recieved a packet") //TODO: remove, will clog up output
			rx_packet := packet_unpack(state.packet_wip[:state.pckt_size])
			// Let the client handler know we got a packet
			rx_chan <- rx_packet
			//reset internal structures
			chunker_reset(state)
		}
	}
	return 0
}

// Owns low level TCP reads
func tcp_socket_read(conn net.Conn, err_chan chan bool, rx_chan chan Packet) {
	// Every connection needs a unique chunker
	// create a zero initilzied state that will be passed to each chunker, this
	// way every connection will get it's own "static" variables
	state := chunker_state{}
	chunker_reset(&state)

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
		chunker(&state, chunk[0:n], n, rx_chan)
	}
err:
	chunker_reset(&state)
	fmt.Println("Closing TCP read socket!")
	conn.Close()
	err_chan <- true
}

// Owns low level TCP writes
func tcp_socket_write(conn net.Conn, err_chan chan bool, shutdown chan bool, write_chan chan Packet) {
	var packet Packet

	for {
		select {
		case packet = <-write_chan:
			break
		case <-shutdown:
			goto err
		}

		packet_binary := packet_pack(packet)

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

// TCP core processed a packet, handle it
// this could be an..
// -> device ack
// -> hello packet (used to set up linux IPC)
// -> login packet (goes to SQL core)
// -> query packet (goes to HTTP core)
func tcp_core_handle_packet_rx(p Packet, cs *Client_state) {

	// Since the device will tell us it's ID,
	// we must wait for the first packet to arive
	// from the devce, we wil use this packet
	// to create the IPC connection to the HTTPs core
	if cs.init != true {
		setup_ipc(p, cs)
		cs.init = true
	}

	//packet_type := p.Packet_type

}
