package main

import (
	"fmt"
	"sync"
	"time"
)

var m map[int16]map_packet // Map to keep track of currently outstanding transactions
var mutex = &sync.Mutex{}

type map_packet struct {
	transaction_id int16
	timestamp      int64
}

func transactions_pop(transaction_id int16) int {
	mutex.Lock()
	// Check to see if a transaction_id is actually in the LL before
	_, ok := m[transaction_id]

	if !ok {
		logger(PRINT_WARN, "Popped a transaction_id %d that was already popped", transaction_id)
		return -1
	}

	mutex.Unlock()
	return int(transaction_id)
}

func transactions_print() {
	for k := range m {
		fmt.Println(m[k])
	}
}

func transaction_scan_timeout(cs *Client_state) {
	mutex.Lock()
	for k := range m {
		if time_ms_since_timestamp(m[k].timestamp) > TCP_PACKET_MS_NO_ACK_CONSIDERED_LOST {
			logger(PRINT_WARN, "Lost packet: ", m[k].transaction_id)

			// create nack
			p := Packet{}
			p.Consumer_ack_req = 0
			p.Packet_type = INTERNAL_ACK_PACKET
			p.Transaction_id = m[k].transaction_id

			// send it to the clien core
			cs.tcp_to_icp_reader_chan <- p

			// pop from map
			delete(m, m[k].transaction_id)
		}
	}
	mutex.Unlock()
}

func transactions_append(transaction_id int16) {
	mutex.Lock()
	// Check if the key is already in the map
	_, ok := m[transaction_id]

	if ok {
		logger(PRINT_FATAL, "FATAL ERROR: Key already preset when trying to add transaction_id:", transaction_id)
	}

	if len(m) >= MAX_OUTSTANDING_TRANSACTIONS {
		logger(PRINT_FATAL, "FATAL ERROR: Outstanding transactions full, unexpected when inserting:", transaction_id)
	}

	m[transaction_id] = map_packet{transaction_id: transaction_id, timestamp: time_ms_since_epoch()}
	mutex.Unlock()
}

func time_ms_since_epoch() int64 {
	return time.Now().UnixNano() / 1e6
}

func time_ms_since_timestamp(timestamp int64) int64 {
	timenow := time.Now().UnixNano() / 1e6
	ret := timenow - timestamp
	if ret < 0 {
		logger(PRINT_FATAL, "Negative timestamp - not sure what happened!")
	}
	return ret
}

func init_transaction_accountant() {
	m = make(map[int16]map_packet)
}

/*
func main() {
	init_transaction_accountant()
	t1 := time_ms_since_epoch()
	time.Sleep(time.Second)
	fmt.Println(time_ms_since_timestamp(t1))

	transactions_append(23)

	time.Sleep(time.Second * 2)
	transaction_scan_timeout()

	transactions_print()
}*/
