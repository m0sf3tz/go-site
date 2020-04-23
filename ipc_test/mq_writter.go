package main

import "bitbucket.org/avd/go-ipc/mq"
import "os"
import "log"

func main() {
	mqd, err := mq.OpenLinuxMessageQueue("smq", os.O_RDWR)
	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, 2)
	b[0] = 0
	b[1] = 231

	mqd.Send(b)
}
