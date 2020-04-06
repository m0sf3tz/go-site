package main

import "fmt"
import "net"
import "io"
import "time"
import "log"

type chunker_state struct {
	packet_parsed_so_far int
	pckt_size            int
	current_parse_type   int
	packet_wip           []byte
}

type client_state struct {
	init                   bool
	sync                   chan string
	ipc_to_tcp_writer_chan chan Packet
	tcp_to_icp_reader_chan chan Packet
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
			fmt.Println(rx_packet)
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

// Handles incoming requests.
func Client_handler(conn net.Conn) {

	// Set up the client state
	cs := client_state{}
	fmt.Println(cs)
	// Set up the timeouts
	set_time_outs(&conn)

	// Initilize the IPC sockets

	// Must create the error channel we will share with the two TCP
	// reader and writter tasks. If any errors occur during read/write/
	// chunking the sub tasks will notify the client hanlder through
	// this error channel

	// the two TCP writer/reader slaves will
	// write into this channel to let client_handler
	// know something is wrong. The tcp_write_shutdown channel
	// is used to let the tcp_socket_writter goroutine to know
	// it is time to shutdown
	err_chan := make(chan bool, 2)
	tcp_write_shutdown := make(chan bool, 1)

	// Handles communication with tcp writer/reader tasks
	tcp_socket_writer_chan := make(chan Packet, MAX_OUTSTANDING_TCP_CORE_SEND)
	tcp_socket_reader_chan := make(chan Packet, MAX_OUTSTANDING_TCP_CORE_SEND)

	cs.ipc_to_tcp_writer_chan = make(chan Packet, MAX_OUTSTANDING_TCP_CORE_SEND) // IPC ---> TCP
	cs.tcp_to_icp_reader_chan = make(chan Packet, MAX_OUTSTANDING_TCP_CORE_SEND) // IPC <--- TCP
	cs.sync = make(chan string, 1)

	// Start the listener and writter goroutines
	go tcp_socket_read(conn, err_chan, tcp_socket_reader_chan)
	go tcp_socket_write(conn, err_chan, tcp_write_shutdown, tcp_socket_writer_chan)

	//
	go ipc_wrangler(cs)

	select {
	case <-err_chan:
		fmt.Println("a TCP error messge was recieved")
		tcp_write_shutdown <- true
		time.Sleep(time.Second * 5)
		goto shutdown_client
		//	case tcp_rx := <-tcp_socket_reader_chan:
		//		tcp_core_handle_packet_rx(tcp_rx, cs)
		//		break
		//case ipc_rx <- ipc_to_tcp_writter_chan:

	}

	// here we handle all todos related to shutting down a client
shutdown_client:
	fmt.Println("closing client_handler")
}

// TCP core processed a packet, handle it
// this could be an..
// -> device ack
// -> hello packet (used to set up linux IPC)
// -> login packet (goes to SQL core)
// -> query packet (goes to HTTP core)
func tcp_core_handle_packet_rx(p Packet, cs client_state) {

	// Since the device will tell us it's ID,
	// we must wait for the first packet to arive
	// from the devce, we wil use this packet
	// to create the IPC connection to the HTTPs core
	if cs.init != true {
		tcp_core_settup_ipc(p, cs)
		cs.init = true
	}

}

func tcp_core_settup_ipc(p Packet, cs client_state) {
	if p.Packet_type != HELLO_WORLD_PACKET {
		log.Fatal("first recieved packet was not a hello-world packet!")
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
	cs.init = true
}
