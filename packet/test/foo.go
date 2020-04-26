package main

import ()

type Packet struct {
	c chan bool
}

func main() {
	p := Packet{}
	p.c = make(chan bool, 1)
	foo(p)
	<-p.c
}
