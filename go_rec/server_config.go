package main

// IF set, will force the ipc socket file to be always == to IPC_TEST_SOCKET_NAME
const IPC_TEST_MODE = true
const SOCKET_PATH = "/tmp/"
const IPC_TEST_SOCKET_NAME = "MISTER_CAT"
const (
	CONN_HOST = "" /* In Go this means listen to ALL the interfaces */
	CONN_PORT = "3334"
	CONN_TYPE = "tcp"
)

// Logging levels
const PRINT_DEBUG = 0
const PRINT_NORMAL = 1
const PRINT_WARN = 2
const PRINT_CRITICAL = 3
const PRINT_FATAL = 4

var CURRENT_LOG_LEVEL int = PRINT_DEBUG

// used to set variables inside the ack/nak subsystem
const MAX_OUTSTANDING_TRANSACTIONS = 16

// Set timeouts

const TCP_PACKET_MS_NO_ACK_CONSIDERED_LOST = 1000
