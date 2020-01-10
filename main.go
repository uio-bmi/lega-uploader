package main

import (
	"./auth"
	"./conf"
	"./resumables"
	"./uploading"
	"flag"
	"fmt"
	"github.com/logrusorgru/aurora"
	"golang.org/x/crypto/ssh/terminal"
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
		conf.Configure()
	}
	if *loginFlag {
		fmt.Println(aurora.Yellow("Username: "))
		var username string
		_, _ = fmt.Scanln(&username)
		fmt.Println(aurora.Yellow("Password: "))
		bytePassword, err := terminal.ReadPassword(0)
		if err != nil {
			log.Fatal(aurora.Red(err))
		}
		password := string(bytePassword)
		if err := auth.Login(username, password); err != nil {
			log.Fatal(aurora.Red(err))
		}
	}
	if *resumablesFlag {
		if err := resumables.Resumables(); err != nil {
			log.Fatal(aurora.Red(err))
		}
	}
	if *uploadFlag {
		args := os.Args[2:]
		for _, file := range args {
			if err := uploading.Upload(file, false); err != nil {
				log.Fatal(aurora.Red(err))
			}
		}
	}
	if *resumeFlag {
		args := os.Args[2:]
		for _, file := range args {
			if err := uploading.Upload(file, true); err != nil {
				log.Fatal(aurora.Red(err))
			}
		}
	}
}
