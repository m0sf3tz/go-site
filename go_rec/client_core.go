package main

import "fmt"
import "time"
import "net"
import "sync"

type Client_state struct {
	init                   bool
	err_chan_tcp           chan bool
	tcp_write_shutdown     chan bool
	tcp_socket_writer_chan chan Packet
	tcp_socket_reader_chan chan Packet

	// Timer -> Client core
	client_event_timer chan bool

	// Unix -> Client core
	ipc_to_tcp_writer_chan chan Packet
	tcp_to_icp_reader_chan chan Packet
	ipc_read_shutdown      chan bool
	ipc_write_shutdown     chan bool
	ipc_shutdown           chan bool

	wg        sync.WaitGroup
	device_id string
	sync      chan bool
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

// Handles incoming requests.
func Client_handler(conn net.Conn) {

	// Set up the client state
	cs := Client_state{}
	// Set up the timeouts
	set_time_outs(&conn)
	// Initilize the packet_accountant
	init_transaction_accountant()

	// Must create the error channel we will share with the two TCP
	// reader and writter tasks. If any errors occur during read/write/
	// chunking the sub tasks will notify the client hanlder through
	// this error channel

	// the two TCP writer/reader slaves will
	// write into this channel to let client_handler
	// know something is wrong. The tcp_write_shutdown channel
	// is used to let the tcp_socket_writter goroutine to know
	// it is time to shutdown
	cs.err_chan_tcp = make(chan bool, 2)
	cs.tcp_write_shutdown = make(chan bool, 1)

	// Handles communication with tcp writer/reader tasks
	cs.tcp_socket_writer_chan = make(chan Packet, MAX_OUTSTANDING_TCP_CORE_SEND)
	cs.tcp_socket_reader_chan = make(chan Packet, MAX_OUTSTANDING_TCP_CORE_SEND)

	cs.client_event_timer = make(chan bool, 1) // Internal to client handler

	cs.ipc_to_tcp_writer_chan = make(chan Packet, MAX_OUTSTANDING_TCP_CORE_SEND) // IPC ---> TCP
	cs.tcp_to_icp_reader_chan = make(chan Packet, MAX_OUTSTANDING_TCP_CORE_SEND) // IPC <--- TCP

	cs.ipc_read_shutdown = make(chan bool, 1)  // IPC (internal)
	cs.ipc_write_shutdown = make(chan bool, 1) // IPC (internal)
	cs.ipc_shutdown = make(chan bool, 1)       // IPC (internal) can write to this to triger IPC shutdown and reset

	cs.sync = make(chan bool, 1)

	go tcp_starter(conn, &cs)
	go ipc_starter(&cs)

	// Will trigger the event timer
	go event_generator(&cs)

	for {
		select {
		case <-cs.err_chan_tcp:
			logger(PRINT_WARN, "A TCP error messge was recieved")
			cs.tcp_write_shutdown <- true
			time.Sleep(time.Second * 5)
			cs.ipc_shutdown <- true
			goto shutdown_client

		case <-cs.ipc_shutdown:
			logger(PRINT_WARN, "shutting down IPC")
			cs.wg.Add(2)
			cs.ipc_read_shutdown <- true
			cs.ipc_write_shutdown <- true
			break

		case tcp_rx := <-cs.tcp_socket_reader_chan:
			logger(PRINT_DEBUG, "sending to ipc output")
			client_core_handle_packet_rx(tcp_rx, &cs)
			break

		case ipc_rx := <-cs.ipc_to_tcp_writer_chan:
			logger(PRINT_DEBUG, "received IPC packet, sending to TCP")
			cs.tcp_socket_writer_chan <- ipc_rx
			client_enqueue_transaction(ipc_rx)
			break

		case <-cs.client_event_timer:
			transaction_scan_timeout(&cs)
			break
		}

	}
	// here we handle all todos related to shutting down a client
shutdown_client:
	// wait for the IPC to shut down
	if waitTimeout(&(cs.wg), time.Second) {
		fmt.Println("Timed out waiting for wait group")
	} else {
		fmt.Println("Wait group finished")
	}

	//TODO: Close IPC and TCP (MUST DO!)
	logger(PRINT_WARN, "closing client_handler")
}

// If a packet requires an ACK, we will track it here
func client_enqueue_transaction(p Packet) {
	if p.Consumer_ack_req != CONSUMER_ACK_REQUIRED {
		return
	}
	logger(PRINT_DEBUG, "Transaction_id ", p.Transaction_id, " needs a device ACK, adding it to the TX map")
	transactions_append(p.Transaction_id)
}

func event_generator(cs *Client_state) {
	for {
		time.Sleep(time.Millisecond * 250)
		cs.client_event_timer <- true
	}
}

// tcp core processed a packet, handle it
// this could be an..
// -> device ack
// -> login packet (goes to SQL core)
// -> query packet (goes to HTTP core)
// TODO: fill up
func client_core_handle_packet_rx(p Packet, cs *Client_state) {
	if p.Consumer_ack_req == CONSUMER_ACK_REQUIRED {
		fmt.Println("Sending ACK to server for transaction_id", p.Transaction_id)
		cs.tcp_socket_writer_chan <- create_ack_pack(p, ACK_GOOD)
	}

	cs.tcp_to_icp_reader_chan <- p
}
