package main

import "os"

//import "bufio"
import "fmt"
import "os/exec"

func clear() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}

var top_menu []string

func init_menu() {
	top_menu = make([]string, 1)
	top_menu = append(top_menu, "1) what's up sam")
}

func main() {

	init_menu()
	/*
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		fmt.Println(text)
	*/

	fmt.Println(top_menu)

}
