// Package main is the main package of lega-uploader command-line tool, containing "config", "login", "upload",
// "resumables" and "resume" commands implementations along with additional helper methods.
package main

import (
	"flag"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/uio-bmi/lega-uploader/auth"
	"github.com/uio-bmi/lega-uploader/conf"
	"github.com/uio-bmi/lega-uploader/resuming"
	"github.com/uio-bmi/lega-uploader/uploading"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	version = "dev"
	date    = "unknown"
)

func main() {
	configFlag := flag.Bool("config", false, "Configure the client.")
	loginFlag := flag.Bool("login", false, "Log in.")
	uploadFlag := flag.Bool("upload", false, "<file1|folder1> <file2|folder2> ... <fileN|folderN>\tUpload files or directories.")
	resumablesFlag := flag.Bool("resumables", false, "List unfinished resumable uploads.")
	resumeFlag := flag.Bool("resume", false, "<file1|folder1> <file2|folder2> ... <fileN|folderN>\tResume files or directories upload.")
	versionFlag := flag.Bool("version", false, "Print tool version.")

	flag.Parse()

	if *configFlag {
		fmt.Println(aurora.Yellow("Instance URL: "))
		var instanceURL string
		_, err := fmt.Scanln(&instanceURL)
		if err != nil {
			log.Fatal(aurora.Red(err))
		}
		if strings.HasPrefix(instanceURL, "http://") {
			log.Fatal(aurora.Red("http protocol is not supported, please use https"))
		}
		if !strings.HasPrefix(instanceURL, "https://") {
			instanceURL = "https://" + instanceURL
		}
		instanceURL = strings.TrimRight(instanceURL, "/")

		configurationProvider, err := conf.NewConfigurationProvider(nil)
		if err != nil {
			log.Fatal(aurora.Red(err))
		}
		configuration := conf.NewConfiguration(instanceURL, nil)
		if err := configurationProvider.SaveConfiguration(configuration); err != nil {
			log.Fatal(aurora.Red(err))
		}
	} else if *loginFlag {
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
	} else if *resumablesFlag {
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
	} else if *uploadFlag {
		uploader, err := uploading.NewUploader(nil, nil, nil)
		if err != nil {
			log.Fatal(aurora.Red(err))
		}
		args := os.Args[2:]
		for _, file := range args {
			if err := uploader.Upload(file, false); err != nil {
				log.Fatal(aurora.Red(err))
			}
		}
	} else if *resumeFlag {
		uploader, err := uploading.NewUploader(nil, nil, nil)
		if err != nil {
			log.Fatal(aurora.Red(err))
		}
		args := os.Args[2:]
		for _, file := range args {
			if err := uploader.Upload(file, true); err != nil {
				log.Fatal(aurora.Red(err))
			}
		}
	} else if *versionFlag {
		fmt.Println(aurora.Blue(version))
		fmt.Println(aurora.Yellow(date))
	} else {
		flag.Usage()
	}
}
