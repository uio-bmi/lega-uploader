package main

import (
	"./auth"
	"./conf"
	"./resuming"
	"./uploading"
	"flag"
	"fmt"
	"github.com/logrusorgru/aurora"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	configFlag := flag.Bool("config", false, "Configure the client.")
	loginFlag := flag.Bool("login", false, "Log in.")
	uploadFlag := flag.Bool("upload", false, "Upload files or directories.")
	resumablesFlag := flag.Bool("resumables", false, "List unfinished resumable uploads.")
	resumeFlag := flag.Bool("resume", false, "Resume files or directories upload.")

	flag.Parse()

	if *configFlag {
		fmt.Println(aurora.Yellow("Instance URL: "))
		var instanceURL string
		_, err := fmt.Scanln(&instanceURL)
		if err != nil {
			log.Fatal(aurora.Red(err))
		}
		configurationProvider, err := conf.NewConfigurationProvider(nil)
		if err != nil {
			log.Fatal(aurora.Red(err))
		}
		configuration := conf.NewConfiguration(strings.TrimRight(instanceURL, "/"), nil)
		if err := configurationProvider.SaveConfiguration(configuration); err != nil {
			log.Fatal(aurora.Red(err))
		}
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
		authenticationManager, err := auth.NewAuthenticationManager(nil, nil)
		if err != nil {
			log.Fatal(aurora.Red(err))
		}
		if err := authenticationManager.Authenticate(username, password); err != nil {
			log.Fatal(aurora.Red(err))
		} else {
			fmt.Println(aurora.Green("Success!"))
		}
	}

	if *resumablesFlag {
		resumablesManager, err := resuming.NewResumablesManager(nil, nil)
		if err != nil {
			log.Fatal(aurora.Red(err))
		}
		resumables, err := resumablesManager.GetResumables()
		if err != nil {
			log.Fatal(aurora.Red(err))
		}
		for _, resumable := range *resumables {
			fmt.Println(aurora.Blue(resumable.Name + "\t (" + strconv.FormatInt(resumable.Size, 10) + " bytes uploaded)"))
		}
	}

	uploader, err := uploading.NewUploader(nil, nil, nil)
	if err != nil {
		log.Fatal(aurora.Red(err))
	}

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
