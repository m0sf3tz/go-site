package main

import "fmt"
import "log"
import "os"
import "bitbucket.org/avd/go-ipc/mq"

func main() {

	dmq, err := mq.OpenLinuxMessageQueue("smq_to_core", os.O_RDWR)
	if err != nil {
		logger(PRINT_FATAL, "Could not create linux smq - (did yu increase the limit in /proc/sys/fs/mqueue/msg_max?  error: ", err)
	}

	b := make([]byte, 1024)
	for {
		n, err := dmq.Receive(b)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(b[:n])
	}
}

//write, readFrom work
