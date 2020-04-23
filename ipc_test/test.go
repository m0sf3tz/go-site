package main

import "bitbucket.org/avd/go-ipc/mq"
import "fmt"
import "log"
import "time"
import "os"
import "sync"

var client_map map[uint32]chan int
var dmq, dmq_2 *mq.LinuxMessageQueue
var mutex_map, mutex_mq sync.Mutex

func register_client(id uint32, c chan int) {
	mutex_map.Lock()
	if c == nil {
		log.Fatal("failed, map was not init")
	}

	if _, ok := client_map[id]; ok {
		log.Fatal("failed, device id already registered: ", id)
	}
	client_map[id] = c
	mutex_map.Unlock()
}

func client_deregister(id uint32) {
	mutex_map.Lock()
	delete(client_map, id)
	mutex_map.Unlock()
}

func send_to_client(id uint32, m int) {
	mutex_map.Lock()

	if _, ok := client_map[id]; !ok {
		log.Fatal("failed, client was not registered: ", id)
	}
	client_map[id] <- m

	mutex_map.Unlock()
}

func init_map() {
	client_map = make(map[uint32]chan int)
}

func init_lmq() {
	mq.DestroyLinuxMessageQueue("smq")
	mq.DestroyLinuxMessageQueue("smq-2")

	var err error
	dmq, err = mq.CreateLinuxMessageQueue("smq", os.O_RDWR, 0666, 4, 512)
	if err != nil {
		log.Fatal(err)
	}

	dmq_2, err = mq.CreateLinuxMessageQueue("smq-2", os.O_RDWR, 0666, 4, 512)
	if err != nil {
		log.Fatal(err)
	}

}

func pa() {
	c := make(chan int)
	register_client(0, c)

	for {
		i := <-c
		fmt.Println("process a got: ", i)
	}

}

func pb() {
	c := make(chan int)
	register_client(1, c)

	for {
		i := <-c
		fmt.Println("process b got: ", i)
	}
}

func pc() {
	c := make(chan int)
	register_client(2, c)

	for {
		i := <-c
		fmt.Println("process c got: ", i)
	}
}

func mq_write(msg []byte) {
	mutex_mq.Lock()
	err := dmq.Send(msg)
	if err != nil {
		log.Fatal("Error writting to Q", err)
	}
	mutex_mq.Unlock()
}

func mq_listner() {
	rx := make([]byte, 512)
	for {
		_, err := dmq.Receive(rx)
		if err != nil {
			log.Fatal("failed to read q")
		}

		var id uint32 = uint32(rx[0])
		var msg int = int(rx[1])

		fmt.Println("to : ", id, " message : ", msg)

		send_to_client(id, msg)
	}
}

func closer() {
	dmq.Destroy()
	fmt.Println("destroyed mq")
}

func main() {
	init_lmq()
	init_map()
	defer closer()
	//go mq_listner()

	go pa()
	go pb()
	go pc()

	time.Sleep(time.Second * 100)

	mq_write([]byte("Sammy"))

}

/*
func main() {
	mq.DestroyLinuxMessageQueue("smq")

	mqd, err := mq.CreateLinuxMessageQueue("smq", os.O_RDWR, 0666, 4, 512)
	if err != nil {
		log.Fatal(err)
	}

}

	mqd.Send([]byte("Sup sam"))

	foo := make([]byte, 20)
	mqd.Receive(foo)

	fmt.Println(string(foo))
}*/
