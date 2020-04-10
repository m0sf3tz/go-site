package main

import "net"
import "fmt"
import "time"

func main() {
	// connect to test IPC socket

	conn, err := net.ListenUnixgram("unixgram", &net.UnixAddr{"/tmp/unixdomain", "unixgram"})
	if err != nil {
		panic(err)
	}

	fmt.Println("here")

	var buf [1024]byte

	conn.WriteToUnix(buf[:], &net.UnixAddr{"/tmp/unixdomain", "unixgram"})
	time.Sleep(time.Second * 5)
}
