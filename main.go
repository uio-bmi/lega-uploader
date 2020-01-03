package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	initialize()

	loginFlag := flag.Bool("login", false, "Log in.")
	uploadFlag := flag.Bool("upload", false, "Upload files or directories.")
	resumablesFlag := flag.Bool("resumables", false, "List unfinished resumable uploads.")
	flag.Parse()
	if *loginFlag {
		login()
	}
	if *resumablesFlag {
		resumables()
	}
	if *uploadFlag {
		args := os.Args[2:]
		for _, file := range args {
			if err := upload(file); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func initialize() {
	if _, err := loadConfiguration(); err != nil {
		createConfiguration("https://localhost")
	}
}
