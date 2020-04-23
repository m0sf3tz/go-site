package main

import (
	"bitbucket.org/avd/go-ipc/mq"
	"log"
	"os"
)

var dmq_from_site_to_core, dmq_from_core_to_site *mq.LinuxMessageQueue
var dmq_from_packet_to_core, dmq_from_core_to_packet *mq.LinuxMessageQueue

func init_lmq_core() {
	mq.DestroyLinuxMessageQueue("smq_site_to_core")
	mq.DestroyLinuxMessageQueue("smq_core_to_site")

	var err error

	//dmq_from_site_to_core, err = mq.CreateLinuxMessageQueue("smq_from_site_to_core", os.O_RDWR, IPC_QUEUE_PERM, IPC_QUEUE_DEPTH, PACKET_LEN_MAX)
	dmq_from_site_to_core, err = mq.CreateLinuxMessageQueue("smq_from_site_to_core", os.O_RDWR, IPC_QUEUE_PERM, IPC_QUEUE_DEPTH, PACKET_LEN_MAX)
	if err != nil {
		logger(PRINT_FATAL, "Could not create linux smq_client - (did yu increase the limit in /proc/sys/fs/mqueue/msg_max?  error: ", err)
	}

	dmq_from_core_to_site, err = mq.CreateLinuxMessageQueue("smq_from_core_to_site", os.O_RDWR, IPC_QUEUE_PERM, IPC_QUEUE_DEPTH, PACKET_LEN_MAX)
	if err != nil {
		logger(PRINT_FATAL, "Could not create linux smq_client - (did yu increase the limit in /proc/sys/fs/mqueue/msg_max?  error: ", err)
	}

	// The rest of the message queues are created by the packet-core, we just need to open them
	dmq_from_packet_to_core, err = mq.OpenLinuxMessageQueue("smq_from_packet_to_core", os.O_RDWR)
	if err != nil {
		logger(PRINT_FATAL, "Could not open linux smq_client", err)
	}

	dmq_from_core_to_packet, err = mq.OpenLinuxMessageQueue("smq_from_core_to_packet", os.O_RDWR)
	if err != nil {
		logger(PRINT_FATAL, "Could not open linux smq_client", err)
	}

	logger(PRINT_DEBUG, "Created linux IPC")
}

func mq_site_to_packet_writter() {
	for {
		rx := make([]byte, PACKET_LEN_MAX)
		n, err := dmq_from_site_to_core.Receive(rx)
		if err != nil {
			log.Fatal("failed to read q")
		}

		dmq_from_core_to_packet.Send(rx[:n])
	}
}

func mq_from_packet_to_site() {
	for {
		rx := make([]byte, PACKET_LEN_MAX)
		n, err := dmq_from_packet_to_core.Receive(rx)
		if err != nil {
			log.Fatal("failed to read q")
		}

		dmq_from_core_to_site.Send(rx[:n])
	}
}

func main() {
	init_lmq_core()

}
