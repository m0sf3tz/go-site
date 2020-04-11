package main

import (
	"bitbucket.org/avd/go-ipc/mq"
	"fmt"
	"log"
	"os"
	"sync"
)

var client_id uint32
var client_map map[uint32]chan Ipc_packet
var dmq_to_core, dmq_to_client *mq.LinuxMessageQueue
var ipc_id, mutex_map, mutex_mq sync.Mutex

func mq_write(msg []byte) error {
	mutex_mq.Lock()
	err := dmq_to_core.Send(msg)
	if err != nil {
		logger(PRINT_FATAL, "Wrote to a dead IPC channel - OR IPC channel full -  should be very rare - restart EVERYTHING", err)
	}
	mutex_mq.Unlock()
	return err
}

func mq_closer() {
	dmq_to_core.Destroy()
	logger(PRINT_WARN, "Closing MQ IPC")
}

func mq_listner() {
	init_map()
	init_lmq()
	defer mq_closer()

	for {
		rx := make([]byte, PACKET_LEN_MAX)
		_, err := dmq_to_client.Receive(rx)
		if err != nil {
			log.Fatal("failed to read q")
		}

		fmt.Println("here")
		ip := ipc_packet_unpack(rx)
		logger(PRINT_DEBUG, "Recieved packed_id :", ip.Id)

		send_to_client(ip)
	}
}

func register_client(id uint32, c chan Ipc_packet) {
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
	// Close the Channel
	close(client_map[id])
	delete(client_map, id)
	mutex_map.Unlock()
}

func send_to_client(ip Ipc_packet) {
	mutex_map.Lock()

	if _, ok := client_map[ip.Id]; !ok {
		log.Fatal("failed, client was not registered: ", ip.Id)
	}
	client_map[ip.Id] <- ip

	mutex_map.Unlock()
}

func init_map() {
	client_map = make(map[uint32]chan Ipc_packet)
}

func init_lmq() {
	mq.DestroyLinuxMessageQueue("smq_to_core")
	mq.DestroyLinuxMessageQueue("smq_to_client")

	var err error

	dmq_to_core, err = mq.CreateLinuxMessageQueue("smq_to_core", os.O_RDWR, IPC_QUEUE_PERM, IPC_QUEUE_DEPTH, PACKET_LEN_MAX)
	if err != nil {
		logger(PRINT_FATAL, "Could not create linux smq_client - (did yu increase the limit in /proc/sys/fs/mqueue/msg_max?  error: ", err)
	}

	dmq_to_client, err = mq.CreateLinuxMessageQueue("smq_to_client", os.O_RDWR, IPC_QUEUE_PERM, IPC_QUEUE_DEPTH, PACKET_LEN_MAX)
	if err != nil {
		logger(PRINT_FATAL, "Could not create linux smq_client - (did yu increase the limit in /proc/sys/fs/mqueue/msg_max?  error: ", err)
	}

	logger(PRINT_DEBUG, "Created linux IPC")
}

func get_ipc_id() uint32 {
	var ret uint32

	ipc_id.Lock()
	ret = client_id
	client_id++
	ipc_id.Unlock()

	return ret
}

func ipc_starter(cs *Client_state) {

	id := get_ipc_id()

	go ipc_reader(id, cs)
	go ipc_writer(id, cs)

}

// Implementation of ipc_writter and reader is
// naive, assumes a lot (chennels will never close....etc)
// tood: fix these?

func ipc_writer(id uint32, cs *Client_state) {
	logger(PRINT_DEBUG, "starting IPC writer")
	var ip Ipc_packet
	ip.Id = id

	for {
		select {
		case ip.P = <-cs.tcp_to_icp_reader_chan:
			break
		case <-cs.ipc_write_shutdown:
			cs.wg.Done()
			return
		}
		logger(PRINT_DEBUG, "ipc_writer sending to linuxQ")

		err := dmq_to_core.Send([]byte("zup sam!"))

		//fmt.Println("fucker", ip.P)
		//	err := mq_write(ipc_packet_pack(ip))
		if err != nil {
			logger(PRINT_FATAL, "Will not handle failure to write to IPC messaeg queue - shut everything down and restart")
		}
	}
}

func ipc_reader(id uint32, cs *Client_state) {
	fmt.Println("starting IPC read server")
	// Register this ipc_reader with the ipc_core, the core will sending
	// any IPC packets to this particular client id into this chan
	var ipc_packet Ipc_packet
	ipc_read_chan := make(chan Ipc_packet, MAX_OUTSTANDING_TRANSACTIONS)
	register_client(id, ipc_read_chan)

	for {
		select {
		case ipc_packet = <-ipc_read_chan:
			break
		case <-cs.ipc_write_shutdown:
			//only once call is required to derigster, will close ipc_read_chan
			client_deregister(id)
			cs.wg.Done()
			return
		}

		logger(PRINT_DEBUG, "IPC_READER read a packet")

		cs.ipc_to_tcp_writer_chan <- ipc_packet.P
	}
}
