# lega-uploader
[![Build Status](https://github.com/uio-bmi/lega-uploader/workflows/Go/badge.svg)](https://github.com/uio-bmi/lega-uploader/actions)
[![GoDoc](https://godoc.org/github.com/uio-bmi/lega-uploader?status.svg)](https://godoc.org/github.com/uio-bmi/lega-uploader)
[![CodeFactor](https://www.codefactor.io/repository/github/uio-bmi/lega-uploader/badge)](https://www.codefactor.io/repository/github/uio-bmi/lega-uploader)
[![Go Report Card](https://goreportcard.com/badge/github.com/uio-bmi/lega-uploader)](https://goreportcard.com/report/github.com/uio-bmi/lega-uploader)
[![codecov](https://codecov.io/gh/uio-bmi/lega-uploader/branch/master/graph/badge.svg)](https://codecov.io/gh/uio-bmi/lega-uploader)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=uio-bmi/lega-uploader)](https://dependabot.com)

[![DeepSource](https://static.deepsource.io/deepsource-badge-light.svg)](https://deepsource.io/gh/uio-bmi/lega-uploader/?ref=repository-badge)

## Installation
To install the latest version of the console app you can use the following one-liner (assuming you are using `bash`):
```
curl -fsSL https://raw.githubusercontent.com/uio-bmi/lega-uploader/master/install.sh | sh
```

Alternatively, go to the [releases page](https://github.com/uio-bmi/lega-uploader/releases) and download the desired binary manually (for example, `.exe` file for Windows).

## Usage
```
$ lega-uploader
Usage of lega-uploader:
  -config
        Configure the client.
  -login
        Log in.
  -resumables
        List unfinished resumable uploads.
  -resume
        <file1|folder1> <file2|folder2> ... <fileN|folderN>     Resume files or directories upload.
  -upload
        <file1|folder1> <file2|folder2> ... <fileN|folderN>     Upload files or directories.
  -version
        Print tool version.
```