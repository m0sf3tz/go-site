package main

import "fmt"

func logger(level int, a ...interface{}) {
	if level >= CURRENT_LOG_LEVEL {
		fmt.Println(a)
	}
	if level == PRINT_FATAL {
		panic(0) // time to die :(
	}
}
