package main

import (
	"fmt"
	"os"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("pfm - Personal Finance CLI Manager")
		fmt.Println("Usage:")
		fmt.Println("  pfm version")
		fmt.Println("  pfm hello")
		return
	}

	switch os.Args[1] {
	case "version":
		fmt.Println("pfm v" + version)
	case "hello":
		fmt.Println("Hello! PFM is working.")
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
	}
}
