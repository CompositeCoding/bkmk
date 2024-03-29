package main

import "fmt"

func log_error(err error, severity int8) {
	if severity == 0 {
		fmt.Printf("INFO: %v", err)
	}
	if severity == 1 {
		fmt.Printf("ERROR: %v", err)
	}
	if severity == 2 {
		fmt.Printf("PANIC: %v", err)
	}
}
