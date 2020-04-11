package main

/**********************************************************
*       Used to define IPC message queue properiets
*********************************************************/

const IPC_QUEUE_DEPTH = (512)
const IPC_QUEUE_SIZE = (PACKET_LEN_MAX)
const IPC_QUEUE_PERM = (0666)

// +4 comes from size of uint32
const MAX_IPC_LEN = (PACKET_LEN_MAX + 4)

type Ipc_packet struct {
	P  Packet
	Id uint32
}
