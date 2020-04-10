package main

import "net"

func main() {
	// connect to test IPC socket

	conn, err := net.Dial("unixgram", "/tmp/unixdomain")
	if err != nil {
		panic(err)
	}

	b := make([]byte, 10)
	conn.Read(b)

	conn.Close()
}

//write, readFrom work
