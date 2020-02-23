package main

import "log"
import "fmt"
import "time"

import "github.com/tarm/serial"

const ACK_SUCCESS = 0x00     // Operation successfully
const ACK_FAIL = 0x01        // Operation failed
const ACK_FULL = 0x04        // Fingerprint database is full
const ACK_NOUSER = 0x05      // No such user
const ACK_USER_EXISTS = 0x07 // already exists
const ACK_TIMEOUT = 0x08     // Acquisition timeout

func checksum(b []byte) byte {
	if len(b) != 8 {
		fmt.Println("Did not get 8 bytes for checksum!")
		panic(0)
	}

	return (b[1] ^ b[2] ^ b[3] ^ b[4] ^ b[5] ^ b[6])
}

func print(b []byte) {
	if len(b) != 8 {
		fmt.Println("Did not get 8 bytes for checksum!")
		panic(0)
	}
	fmt.Printf("0x%X 0x%X 0x%X 0x%X 0x%X 0x%X 0x%X 0x%X \n", b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7])
}

func addUser(ser *serial.Port, id byte) {
	time.Sleep(time.Second * 3)
	res := make([]byte, 8)

	for i := 1; i < 4; i++ {
		command := []byte{0xF5, byte(i), 0, 2, 1, 0, -0x0, 0xf5}
		sum := checksum(command)
		command[6] = sum
		ser.Write(command)
		print(command)

		time.Sleep(time.Second)
		ser.Read(res)
		print(res)
		fmt.Println()
	}

	var check byte = res[4] //this guy contains ack_nouser, timeout, or user priv

	if check == ACK_SUCCESS {
		fmt.Println("Added new user")
	} else if check == ACK_FAIL {
		fmt.Println("Failed to ass new user")
	} else if check == ACK_USER_EXISTS {
		fmt.Println("User already exists")
	} else if check == ACK_TIMEOUT {
		fmt.Println("Timed out adding user")
	} else {
		fmt.Println("not sure what happened...")
	}
}

func totalUsers(ser *serial.Port) {
	command := []byte{0xF5, 0x09, 0x0, 0x0, 00, 0, -0x0, 0xf5}
	sum := checksum(command)
	command[6] = sum
	fmt.Println()
	print(command)
	ser.Write(command)
	res := make([]byte, 8)
	ser.Read(res)
	print(res)
	fmt.Println("Total number of users is: ", res[2]|res[3])
}

func printId(id int) {
	if id == 2 {
		fmt.Println()
		fmt.Println("Hi sam..")
	}
}

func comp(ser *serial.Port) int {
	command := []byte{0xF5, 0x0C, 0x0, 0x0, 00, 0, -0x0, 0xf5}
	sum := checksum(command)
	command[6] = sum
	ser.Write(command)
	res := make([]byte, 8)
	ser.Read(res)

	print(res)
	var check byte = res[4] //this guy contains ack_nouser, timeout, or user priv
	if check == 05 {
		fmt.Println("no user")
		return 0
	}

	fmt.Println("User is: ", res[2]|res[3])
	return int(res[2] | res[3])
}

func eraseAll(ser *serial.Port) {
	command := []byte{0xF5, 0x05, 0x0, 0x0, 00, 0, -0x0, 0xf5}
	sum := checksum(command)
	command[6] = sum
	ser.Write(command)

	res := make([]byte, 8)
	ser.Read(res)
	var check byte = res[4] //this guy contains ack_nouser, timeout, or user priv
	if check == ACK_SUCCESS {
		fmt.Println("Deleted all users")
	} else if check == ACK_FAIL {
		fmt.Println("Failed to delete all users")
	} else {
		fmt.Println("Not sure what happened....")
	}
}

func main() {
	c := &serial.Config{Name: "/dev/ftdi", Baud: 19200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	//totalUsers(s)
	//eraseAll(s)
	id := comp(s)
	printId(id)
	//addUser(s, 0xA5)
}
