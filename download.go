package main

import (
	"encoding/json"
	"fmt"
	"github.com/rs/xid"
	"io"
	"net/http"
	"os"
	"time"
)

//Creates filepath, calls DownloadFile on the filepath and url, and updates download information
func downloadSingleFile(url string, downloadId string, status *string, files map[string]string){

	//Logs url to console window
	fmt.Println(url)

	//Creates filepath to store download using unique id generated and stored in fileId
	fileId := xid.New()
	filePath := "/tmp/" + downloadId + "-" + fileId.String()
	err := DownloadFile(filePath, url)
	if err != nil {
		*status = "FAILURE"
		panic(err)
	} else {
		//Adds file name and url to download information
		files[url] = filePath
	}
}

//DownloadFile will download a url to a local file
func DownloadFile(filepath string, url string) error {

	//Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	//Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func serialDownload(w http.ResponseWriter, r *http.Request, payload Payload){
	status := "PENDING"
	downloadId := xid.New().String()
	files := make(map[string]string)
	startTime := time.Now()
	for _,url := range payload.Urls{
		downloadSingleFile(url, downloadId, &status, files)
	}
	if status == "PENDING" {
		status = "SUCCESSFUL"
	}
	endTime := time.Now()
	DownloadsInfo[downloadId] = DownloadInfo{
		Id :            downloadId,
		Start_time :    startTime,
		End_time :      endTime,
		Status :        status,
		Download_type : payload.Type,
		Files :         files,
	}
	responseid,_ := json.Marshal(Response{Id: downloadId})
	w.Header().Set("Content-Type","application/json")
	w.Write(responseid)
}

func concurrentDownload(w http.ResponseWriter, r *http.Request, payload Payload){
	status := "PENDING"
	downloadId := xid.New().String()
	files := make(map[string]string)
	responseid,_ := json.Marshal(Response{Id: downloadId})
	w.Header().Set("Content-Type","application/json")
	w.Write(responseid)
	var ch = make(chan string)
	for i := 0; i < numberOfLinks; i++ {
		go func() {
			for {
				url, ok := <-ch
				if !ok {
					return
				}
				downloadSingleFile(url, downloadId, &status, files)
			}
		}()
	}
	startTime := time.Now()
	endTime := time.Now()
	DownloadsInfo[downloadId] = DownloadInfo{
		Id :            downloadId,
		Start_time :    startTime,
		End_time :      endTime,
		Status :        status,
		Download_type : payload.Type,
		Files :         files,
	}
	go func() {
		for _, url := range payload.Urls {
			ch <- url
		}
		close(ch)
		if status == "PENDING" {
			status = "SUCCESSFUL"
		}
		endTime = time.Now()
		DownloadsInfo[downloadId] = DownloadInfo{
			Id :            downloadId,
			Start_time :    startTime,
			End_time :      endTime,
			Status :        status,
			Download_type : payload.Type,
			Files :         files,
		}
		return
	}()
}

