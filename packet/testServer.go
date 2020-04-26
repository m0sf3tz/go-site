package main

import "net"
import "time"
import "log"
import "fmt"

func main() {
	// connect to test IPC socket

	fmt.Println("hia")
	_, err := net.Dial("unix", SOCKET_PATH+IPC_TEST_SOCKET_NAME)
	if err != nil {
		log.Fatal(err)
		time.Sleep(time.Second * 100)
	}
}
