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

## Configuration
Before using the app, make sure all the environment variables required for authentication are set:

```
export CENTRAL_EGA_USERNAME=...
export CENTRAL_EGA_PASSWORD=...
export ELIXIR_AAI_TOKEN=...
```

NB: `ELIXIR_AAI_TOKEN` has an expiration time of nearly two hours, so one would need to re-obtain and re-set it upon expiration.

Also, the tool is pre-configured to work with Norwegian Federated EGA instance: https://ega.elixir.no 
If you want to specify another instance, you can set `LOCAL_EGA_INSTANCE_URL` environment variable. 

## Usage

```
$ lega-uploader
lega-uploader [files | resumables | upload] <args>

 files:
  -l, --list    Lists uploaded files
  -d, --delete= Deletes uploaded file by name

 resumables:
  -l, --list    Lists resumable uploads
  -d, --delete= Deletes resumable upload by ID

 upload:
  -f, --file=FILE    File or folder to upload
  -r, --resume       Resumes interrupted upload
```
