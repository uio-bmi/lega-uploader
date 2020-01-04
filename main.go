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
	resumeFlag := flag.Bool("resume", false, "Resume files or directories upload.")
	flag.Parse()
	if *loginFlag {
		login()
	}
	if *resumablesFlag {
		resumables()
	}
	var function func(path string) error
	if *uploadFlag {
		function = upload
	} else {
		function = resume
	}
	if *uploadFlag || *resumeFlag {
		args := os.Args[2:]
		for _, file := range args {
			if err := function(file); err != nil {
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
