// Package main is the main package of lega-uploader command-line tool, containing "upload", "resumables" and "resume"
// commands implementations along with additional helper methods.
package main

import (
	"flag"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/uio-bmi/lega-uploader/resuming"
	"github.com/uio-bmi/lega-uploader/uploading"
	"log"
	"os"
	"strconv"
)

var (
	version = "dev"
	date    = "unknown"
)

func main() {
	uploadFlag := flag.Bool("upload", false, "<file1|folder1> <file2|folder2> ... <fileN|folderN>\tUpload files or directories.")
	resumablesFlag := flag.Bool("resumables", false, "List unfinished resumable uploads.")
	resumeFlag := flag.Bool("resume", false, "<file1|folder1> <file2|folder2> ... <fileN|folderN>\tResume files or directories upload.")
	versionFlag := flag.Bool("version", false, "Print tool version.")

	flag.Parse()

	if *resumablesFlag {
		resumablesManager, err := resuming.NewResumablesManager(nil)
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
		uploader, err := uploading.NewUploader(nil, nil)
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
		uploader, err := uploading.NewUploader(nil, nil)
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
