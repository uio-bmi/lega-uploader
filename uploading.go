package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func upload(path string) error {
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

func uploadFile(file *os.File, stat os.FileInfo) error {
	configuration, err := loadConfiguration()
	if err != nil {
		return err
	}

	fileName := filepath.Base(file.Name())
	var uploadId string

	buffer := make([]byte, *configuration.ChunkSize)
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
		checksum := hex.EncodeToString(sum[:16])
		request, err := http.NewRequest("PATCH", *configuration.InstanceURL+"/stream/"+fileName, bytes.NewReader(chunk))
		if err != nil {
			return err
		}
		request.Header.Add("Authorization", "Bearer "+*configuration.InstanceToken)
		query := request.URL.Query()
		if i != 1 {
			query.Add("uploadId", uploadId)
		}
		query.Add("chunk", strconv.Itoa(i))
		query.Add("md5", checksum)
		request.URL.RawQuery = query.Encode()
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			return err
		}
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
	}

	hashFunction := md5.New()
	_, err = io.Copy(hashFunction, file)
	if err != nil {
		return err
	}
	checksum := hex.EncodeToString(hashFunction.Sum(nil)[:16])
	request, err := http.NewRequest("PATCH", *configuration.InstanceURL+"/stream/"+fileName, bytes.NewReader(buffer))
	if err != nil {
		return err
	}
	request.Header.Add("Authorization", "Bearer "+*configuration.InstanceToken)
	query := request.URL.Query()
	query.Add("uploadId", uploadId)
	query.Add("chunk", "end")
	query.Add("fileSize", strconv.FormatInt(stat.Size(), 10))
	query.Add("md5", checksum)
	request.URL.RawQuery = query.Encode()
	_, err = http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	return nil
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
