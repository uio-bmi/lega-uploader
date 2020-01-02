package main

import (
	"flag"
)

func main() {
	initialize()

	loginFlag := flag.Bool("login", false, "Logs you in.")
	flag.Parse()
	if *loginFlag {
		login()
	}
}

func initialize() {
	if _, err := loadConfiguration(); err != nil {
		createConfiguration("https://localhost")
	}
}
