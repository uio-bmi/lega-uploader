package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	configFlag := flag.Bool("config", false, "Configure the client.")
	loginFlag := flag.Bool("login", false, "Log in.")
	uploadFlag := flag.Bool("upload", false, "Upload files or directories.")
	resumablesFlag := flag.Bool("resumables", false, "List unfinished resumable uploads.")
	resumeFlag := flag.Bool("resume", false, "Resume files or directories upload.")
	flag.Parse()
	if *configFlag {
		configure()
	}
	if *loginFlag {
		login()
	}
	if *resumablesFlag {
		resumables()
	}
	if *uploadFlag {
		args := os.Args[2:]
		for _, file := range args {
			if err := upload(file, false); err != nil {
				log.Fatal(err)
			}
		}
	}
	if *resumeFlag {
		args := os.Args[2:]
		for _, file := range args {
			if err := upload(file, true); err != nil {
				log.Fatal(err)
			}
		}
	}
}
