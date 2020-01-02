package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	initialize()

	loginFlag := flag.Bool("login", false, "Logs you in.")
	uploadFlag := flag.Bool("upload", false, "Uploads file.")
	flag.Parse()
	if *loginFlag {
		login()
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
