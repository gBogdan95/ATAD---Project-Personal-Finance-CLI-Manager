package main

import (
	"fmt"
	"os"

	"github.com/gBogdan95/ATAD---Project-Personal-Finance-CLI-Manager/internal/db"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "version":
		fmt.Println("pfm v" + version)
	case "hello":
		fmt.Println("Hello! PFM is working ✅")
	case "init":
		if err := cmdInit(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command: %s\n\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("pfm - Personal Finance CLI Manager")
	fmt.Println("Usage:")
	fmt.Println("  pfm version")
	fmt.Println("  pfm hello")
	fmt.Println("  pfm init        # create local sqlite db + tables")
}

func cmdInit() error {
	_, err := db.EnsureDir()
	if err != nil {
		return err
	}

	path, err := db.DBPath()
	if err != nil {
		return err
	}

	conn, err := db.Open(path)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := db.InitSchema(conn); err != nil {
		return err
	}

	fmt.Println("Initialized database at:", path)
	return nil
}
