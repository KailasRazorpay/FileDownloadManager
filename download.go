package main

import (
	"fmt"
	"github.com/rs/xid"
	"io"
	"net/http"
	"os"
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

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
