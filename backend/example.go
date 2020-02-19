package main

import (
	"fmt"
)

func main() {
	foo := []string{"string slice"}
	fmt.Println(foo)
	fmt.Println(len(foo))
	if foo[0] == "string slice" {
		fmt.Println("fuker")
	}

}
