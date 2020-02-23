package main

import "fmt"
import "os"

import "github.com/tarm/serial"
import "os/exec"

//import "io"

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
	fmt.Printf("Command : 0x%X 0x%X 0x%X 0x%X 0x%X 0x%X 0x%X 0x%X \n", b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7])
}

func main() {
	fi, err := os.OpenFile("/dev/ftdi", os.O_RDWR, 0660)
	if err != nil {
		fmt.Println("Could not open, bailing out")
		panic(0)
	}

	cmd := exec.Command("stty", "-F", "/dev/ftdi", "sane")
	//cmd := exec.Command("stty", "-F", "/dev/ftdi", "19200", "cread", "clocal")
	err = cmd.Run()
	if err != nil {
		fmt.Println("could not set ttl")
		panic(0)
	}

	command := []byte{0xF5, 0x2c, 0, 0, 0, 0, 0, 0xf5}
	sum := checksum(command)
	command[6] = sum

	print(command)

	_, err = fi.Write(command)
	if err != nil {
		fmt.Println("Failed to write!")
		panic(0)
	}

	res := make([]byte, 1)

	fmt.Println("here1 ")
	n, err := fi.Read(res)
	fmt.Println("here")

	fmt.Println("Read: ", n)
	if err != nil {
		fmt.Println("Failed to write!")
		panic(0)
	}

	defer fi.Close()
}
