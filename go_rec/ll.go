package main

import (
	"container/list"
	"fmt"
)

const MAX_Q_LEN = 10
const MAX_OUTSTANDING = 8

type ll_type_e int

const RX_LL ll_type_e = 0
const TX_LL ll_type_e = 1

const MARKER_FREE = 0
const MARKER_TAKEN = 1

var marker_status_rx [MAX_Q_LEN]int
var marker_status_tx [MAX_Q_LEN]int

type Message struct {
	Marker int
	data   []byte
}

func get_free_marker(t ll_type_e) int {
	var mptr *[10]int

	if t == RX_LL {
		mptr = &marker_status_rx
	} else {
		mptr = &marker_status_tx
	}

	for i := 0; i < MAX_OUTSTANDING; i++ {
		if mptr[i] == MARKER_FREE {
			mptr[i] = MARKER_TAKEN
			return i
		}
	}

	fmt.Println("ran out of free markers")
	return -1
}

func free_marker(t ll_type_e, marker int) int {
	if marker >= MAX_OUTSTANDING || marker < 0 {
		return -1
	}

	if t == RX_LL {
		marker_status_rx[marker] = MARKER_FREE
	} else {
		marker_status_tx[marker] = MARKER_FREE
	}
	return 0
}

func Linked_remove(i int) int {
	for element := l.Front(); element != nil; element = element.Next() {
		if element.Value.(Message).Marker == i {
			l.Remove(element)
			return 0
		}
	}
	return -1
}

func Linked_append(m Message) {
	l.PushBack(m)
}

func Linked_pop() Message {
	m := l.Back()
	if m == nil {
		fmt.Println("Error... tried to pop an empty list")
		panic(0)
	}
	l.Remove(l.Back())
	return m.Value.(Message)
}

var l *list.List

func main() {
	fmt.Println("Go Linked Lists Tutorial")
	l = list.New()
	Linked_append(Message{1, nil})
	Linked_append(Message{2, nil})
	Linked_append(Message{3, nil})
	// we now have a linked list with '1' at the back of the list
	// and '2' at the front of the list.

	Linked_remove(1)

	for element := l.Front(); element != nil; element = element.Next() {
		// do something with element.Value
		fmt.Println(element.Value)
	}

}
