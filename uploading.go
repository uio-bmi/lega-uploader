package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func upload(path string) error {
	var err error
	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			return err
		}
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return uploadFolder(file)
	} else {
		return uploadFile(file, stat)
	}
}

func uploadFolder(folder *os.File) error {
	readdir, err := folder.Readdir(-1)
	if err != nil {
		return err
	}
	for _, file := range readdir {
		abs, err := filepath.Abs(filepath.Join(folder.Name(), file.Name()))
		if err != nil {
			return err
		}
		err = upload(abs)
		if err != nil {
			return err
		}
	}
	return nil
}

func uploadFile(file *os.File, stat os.FileInfo) error {
	totalSize := stat.Size()
	println("Uploading file: " + file.Name() + " (" + strconv.FormatInt(totalSize, 10) + " bytes)")
	bar := pb.StartNew(100)
	bar.Start()
	configuration, err := loadConfiguration()
	if err != nil {
		return err
	}

	fileName := filepath.Base(file.Name())
	var uploadId string

	buffer := make([]byte, *configuration.ChunkSize*1024*1024)
	for i := 1; true; i++ {
		read, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		chunk := buffer[:read]
		sum := md5.Sum(chunk)
		params := map[string]string{
			"chunk": strconv.Itoa(i),
			"md5":   hex.EncodeToString(sum[:16])}
		if i != 1 {
			params["uploadId"] = uploadId
		}
		response, err := doRequest(http.MethodPatch,
			*configuration.InstanceURL+"/stream/"+fileName,
			bytes.NewReader(chunk),
			map[string]string{"Authorization": "Bearer " + *configuration.InstanceToken},
			params,
			nil,
			nil)
		if err != nil {
			return err
		}
		//noinspection GoDeferInLoop
		defer response.Body.Close()
		if response.StatusCode != 200 {
			return errors.New(response.Status)
		}
		jsonResponse := make(map[string]interface{})
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&jsonResponse)
		if err != nil {
			return err
		}
		statusCode := fmt.Sprint(jsonResponse["statusCode"])
		if statusCode != "201" {
			return errors.New("Status code: " + statusCode)
		}
		uploadId = fmt.Sprint(jsonResponse["id"])
		bar.Add64(int64(read) * 100 / totalSize)
	}
	bar.SetCurrent(100)
	hashFunction := md5.New()
	_, err = io.Copy(hashFunction, file)
	if err != nil {
		return err
	}
	checksum := hex.EncodeToString(hashFunction.Sum(nil)[:16])
	response, err := doRequest(http.MethodPatch,
		*configuration.InstanceURL+"/stream/"+fileName,
		nil,
		map[string]string{"Authorization": "Bearer " + *configuration.InstanceToken},
		map[string]string{"uploadId": uploadId,
			"chunk":    "end",
			"fileSize": strconv.FormatInt(totalSize, 10),
			"md5":      checksum},
		nil,
		nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New(response.Status)
	}
	bar.Finish()
	return nil
}
