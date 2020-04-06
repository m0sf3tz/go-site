package main

import "fmt"
import "net"
import "time"

type Client_state struct {
	init                   bool
	sync                   chan string
	ipc_to_tcp_writer_chan chan Packet
	tcp_to_icp_reader_chan chan Packet
	client_event_timer     chan bool
}

// Handles incoming requests.
func Client_handler(conn net.Conn) {

	// Set up the client state
	cs := Client_state{}
	fmt.Println(cs)
	// Set up the timeouts
	set_time_outs(&conn)
	// Initilize the packet_accountant
	init_transaction_accountant()
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

	cs.client_event_timer = make(chan bool, 1)                                   // Internal to client handler
	cs.ipc_to_tcp_writer_chan = make(chan Packet, MAX_OUTSTANDING_TCP_CORE_SEND) // IPC ---> TCP
	cs.tcp_to_icp_reader_chan = make(chan Packet, MAX_OUTSTANDING_TCP_CORE_SEND) // IPC <--- TCP
	cs.sync = make(chan string, 1)

	// Will trigger the event timer
	go event_generator(&cs)

	// Start the listener and writter goroutines
	go tcp_socket_read(conn, err_chan, tcp_socket_reader_chan)
	go tcp_socket_write(conn, err_chan, tcp_write_shutdown, tcp_socket_writer_chan)

	go ipc_wrangler(&cs)

	for {
		select {
		case <-err_chan:
			fmt.Println("a TCP error messge was recieved")
			tcp_write_shutdown <- true
			time.Sleep(time.Second * 5)
			goto shutdown_client
		case tcp_rx := <-tcp_socket_reader_chan:
			tcp_core_handle_packet_rx(tcp_rx, &cs)
			break
		case ipc_rx := <-cs.ipc_to_tcp_writer_chan:
			fmt.Println("sending to tcp")
			client_enqueue_transaction(ipc_rx)
			break
		case <-cs.client_event_timer:
			fmt.Println("here??")
			transaction_scan_timeout(&cs)
			break
		}

	}
	// here we handle all todos related to shutting down a client
shutdown_client:
	//TODO: erase socket
	fmt.Println("closing client_handler")
}

// If a packet requires an ACK, we will track it here
func client_enqueue_transaction(p Packet) {
	if p.Consumer_ack_req != CONSUMER_ACK_REQUIRED {
		return
	}
	transactions_append(p.Transaction_id)
	transactions_print()
}

func event_generator(cs *Client_state) {
	for {
		time.Sleep(time.Millisecond * 250)
		cs.client_event_timer <- true
	}
}
