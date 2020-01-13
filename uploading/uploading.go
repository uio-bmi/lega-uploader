package uploading

import (
	"../conf"
	"../requests"
	"../resumables"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/cheggaaa/pb/v3"
	"github.com/logrusorgru/aurora"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

func Upload(path string, resume bool) error {
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
		return uploadFolder(file, resume)
	} else {
		if resume {
			fileName := filepath.Base(file.Name())
			resumablesList, err := resumables.GetResumables()
			if err != nil {
				return err
			}
			for _, resumable := range *resumablesList {
				if resumable.Name == fileName {
					return uploadFile(file, stat, &resumable.Id, resumable.Size, resumable.Chunk)
				}
			}
			return nil
		} else {
			return uploadFile(file, stat, nil, 0, 1)
		}
	}
}

func uploadFolder(folder *os.File, resume bool) error {
	readdir, err := folder.Readdir(-1)
	if err != nil {
		return err
	}
	for _, file := range readdir {
		abs, err := filepath.Abs(filepath.Join(folder.Name(), file.Name()))
		if err != nil {
			return err
		}
		err = Upload(abs, resume)
		if err != nil {
			return err
		}
	}
	return nil
}

func uploadFile(file *os.File, stat os.FileInfo, uploadId *string, offset int64, startChunk int64) error {
	totalSize := stat.Size()
	fmt.Println(aurora.Blue("Uploading file: " + file.Name() + " (" + strconv.FormatInt(totalSize, 10) + " bytes)"))
	bar := pb.StartNew(100)
	bar.SetCurrent(offset * 100 / totalSize)
	bar.Start()
	configuration, err := conf.LoadConfiguration()
	if err != nil {
		return err
	}

	fileName := filepath.Base(file.Name())

	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}
	buffer := make([]byte, *configuration.ChunkSize*1024*1024)
	client := requests.NewClient()
	for i := startChunk; true; i++ {
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
			"chunk": strconv.FormatInt(i, 10),
			"md5":   hex.EncodeToString(sum[:16])}
		if i != 1 {
			params["uploadId"] = *uploadId
		}
		response, err := client.DoRequest(http.MethodPatch,
			*configuration.InstanceURL+"/stream/"+url.QueryEscape(fileName),
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
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		if uploadId == nil {
			uploadId = new(string)
		}
		*uploadId, err = jsonparser.GetString(body, "id")
		if err != nil {
			return err
		}
		bar.Add64(int64(read) * 100 / totalSize)
	}
	bar.SetCurrent(100)
	hashFunction := md5.New()
	_, err = io.Copy(hashFunction, file)
	if err != nil {
		return err
	}
	checksum := hex.EncodeToString(hashFunction.Sum(nil)[:16])
	response, err := client.DoRequest(http.MethodPatch,
		*configuration.InstanceURL+"/stream/"+url.QueryEscape(fileName),
		nil,
		map[string]string{"Authorization": "Bearer " + *configuration.InstanceToken},
		map[string]string{"uploadId": *uploadId,
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