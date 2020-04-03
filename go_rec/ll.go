package main

import (
	"fmt"
	"log"
	"time"
)

const MAX_OUTSTANDING_TRANSACTIONS = 16

var m map[int16]Packet // Map to keep track of currently outstanding transactions

type Packet struct {
	transaction_id int16
	timestamp      int64
}

func Transactions_pop(transaction_id int16) int {
	// Check to see if a transaction_id is actually in the LL before
	_, ok := m[transaction_id]

	if !ok {
		log.Printf("Popped a transaction_id %d that was already popped", transaction_id)
		return -1
	}

	delete(m, transaction_id)

	return int(transaction_id)
}

func Transactions_print() {
	for k := range m {
		fmt.Println(m[k])
	}
}

func Transactions_append(transaction_id int16) {
	// Check if the key is already in the map
	_, ok := m[transaction_id]

	if ok {
		log.Fatal("FATAL ERROR: Key already preset when trying to add transaction_id:", transaction_id)
	}

	if len(m) >= MAX_OUTSTANDING_TRANSACTIONS {
		log.Fatal("FATAL ERROR: Outstanding transactions full, unexpected when inserting:", transaction_id)
	}

	m[transaction_id] = Packet{transaction_id: transaction_id, timestamp: time_ms_since_epoch()}
}

func time_ms_since_epoch() int64 {
	return time.Now().UnixNano() / 1e6
}

func main() {
	m = make(map[int16]Packet)

	var i int16
	for ; i < 20; i++ {
		Transactions_append(i)
	}

	/*
		m := make(map[int]Packet)
		m[23] = Packet{transaction_id: 23, timestamp: time_ms_since_epoch()}
		_, ok := m[232]
		fmt.Println(ok)
		delete(m, 23)
		fmt.Println(m)
		/*
			l = list.New()
			Linked_append(23)
			Linked_append(12)
			Linked_append(67)
			Linked_print()
	*/
	/*
		fmt.Println("Go Linked Lists Tutorial")
		l = list.New()
		Linked_append(Packet{1, nil})
		Linked_append(Packet{2, nil})
		Linked_append(Packet{3, nil})
		// we now have a linked list with '1' at the back of the list
		// and '2' at the front of the list.

		Linked_remove(1)

		for element := l.Front(); element != nil; element = element.Next() {
			// do something with element.Value
			fmt.Println(element.Value)
		}
	*/
}
