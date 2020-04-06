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
