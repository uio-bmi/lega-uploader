# lega-uploader
[![Build Status](https://github.com/uio-bmi/lega-uploader/workflows/Go/badge.svg)](https://github.com/uio-bmi/lega-uploader/actions)
[![GoDoc](https://godoc.org/github.com/uio-bmi/lega-uploader?status.svg)](https://pkg.go.dev/github.com/uio-bmi/lega-uploader?tab=subdirectories)
[![CodeFactor](https://www.codefactor.io/repository/github/uio-bmi/lega-uploader/badge)](https://www.codefactor.io/repository/github/uio-bmi/lega-uploader)
[![Go Report Card](https://goreportcard.com/badge/github.com/uio-bmi/lega-uploader)](https://goreportcard.com/report/github.com/uio-bmi/lega-uploader)
[![codecov](https://codecov.io/gh/uio-bmi/lega-uploader/branch/master/graph/badge.svg)](https://codecov.io/gh/uio-bmi/lega-uploader)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=uio-bmi/lega-uploader)](https://dependabot.com)

[![DeepSource](https://static.deepsource.io/deepsource-badge-light.svg)](https://deepsource.io/gh/uio-bmi/lega-uploader/?ref=repository-badge)

## Installation

### Linux
```
curl -fsSL https://raw.githubusercontent.com/uio-bmi/lega-uploader/master/install.sh | sudo sh
```

### MacOS
```
curl -fsSL https://raw.githubusercontent.com/uio-bmi/lega-uploader/master/install.sh | sh
```

### Windows
Go to the [releases page](https://github.com/uio-bmi/lega-uploader/releases) and download the binary manually.

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
