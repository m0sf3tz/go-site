package main

import "os"

import "bufio"
import "fmt"
import "os/exec"

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)

	exec.Command("clear") //Linux example, its tested
}
