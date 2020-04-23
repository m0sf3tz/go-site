package main

import "bitbucket.org/avd/go-ipc/mq"
import "os"
import "log"
import "fmt"

func main() {
	mqd, err := mq.OpenLinuxMessageQueue("smq", os.O_RDWR)
	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, 20)

	n, err := mqd.Receive(b)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b[:n]))

	for {
	}
}
