# lega-uploader
[![Build Status](https://github.com/uio-bmi/lega-uploader/workflows/Go/badge.svg)](https://github.com/uio-bmi/lega-uploader/actions)
[![CodeFactor](https://www.codefactor.io/repository/github/uio-bmi/lega-uploader/badge)](https://www.codefactor.io/repository/github/uio-bmi/lega-uploader)

[![DeepSource](https://static.deepsource.io/deepsource-badge-light.svg)](https://deepsource.io/gh/uio-bmi/lega-uploader/?ref=repository-badge)

## Installation
```
sudo curl -L "https://github.com/uio-bmi/lega-uploader/releases/download/v0.0.4/lega-uploader_$(uname -s)_$(uname -m)" -o /usr/local/bin/lega-uploader && sudo chmod +x /usr/local/bin/lega-uploader
```

## Usage
```
Usage of lega-uploader:
  -config
    	Configure the client.
  -login
    	Log in.
  -resumables
    	List unfinished resumable uploads.
  -resume
    	Resume files or directories upload.
  -upload
    	Upload files or directories.
```