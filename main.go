package main

import (
	"./auth"
	"./conf"
	"./resuming"
	"./uploading"
	"flag"
	"github.com/logrusorgru/aurora"
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
		auth.Login()
	}
	if *resumablesFlag {
		resuming.Resumables()
	}
	uploader := uploading.NewUploader(nil, nil, nil)
	if *uploadFlag {
		args := os.Args[2:]
		for _, file := range args {
			if err := uploader.Upload(file, false); err != nil {
				log.Fatal(aurora.Red(err))
			}
		}
	}
	if *resumeFlag {
		args := os.Args[2:]
		for _, file := range args {
			if err := uploader.Upload(file, true); err != nil {
				log.Fatal(aurora.Red(err))
			}
		}
	}
}
