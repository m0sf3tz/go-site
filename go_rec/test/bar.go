package main

import "time"

func foo(p Packet) {
	time.Sleep(time.Second)
	p.c <- true
}
