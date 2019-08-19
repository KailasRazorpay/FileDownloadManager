package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"github.com/rs/xid"
	"log"
	"net/http"
	"strings"
)

type DownloadInfo struct{
	Id string
	Start_time string
	End_time string
	Status string
	Download_type string
	Files map[string]string
}

type Payload struct{
	Type string
	Urls []string
}

type Response struct{
	Id string
}

var DownloadsInfo = map[string]DownloadInfo{}

type HealthHandler struct{}

func (re HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(200)
	fmt.Fprintln(w,"You have hit the health tag")
}

type StatusHandler struct{}

func (st StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	segments := strings.Split(r.URL.Path, "/")
	id := segments[2]
	fmt.Println(DownloadsInfo[id])
	responseBody,err := json.Marshal(DownloadsInfo[id])
	if err != nil{
		panic(err)
	}
	w.Header().Set("Content-Type","application/json")
	w.Write(responseBody)
}

func downloadSingleFile(url string, download_id string, status *string, files map[string]string){
	fmt.Println(url)
	file_id := xid.New()
	filepath := "/tmp/"+download_id+"-"+file_id.String()
	err := DownloadFile(filepath, url)
	if err != nil {
		*status = "FAILURE"
		panic(err)
	} else {
		files[url] = filepath
	}
}

type DownloadHandler struct{}

func (d DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	payload := Payload{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		panic(err)
	}
	if(payload.Type == "serial"){
		status := "SUCCESSFUL"
		download_id := xid.New().String()
		files := make(map[string]string)
		DownloadsInfo[download_id] = DownloadInfo{
			Id : download_id,
			Start_time : "2019-08-18 09:46:42.89588 +0530 IST m=+12.135120014",
			//End_time : "2019-08-18 09:48:46.287554 +0530 IST m=+135.526087303",
			//Status : status,
			Download_type : "serial",
			Files : files,
		}
		for _,url := range payload.Urls{
			downloadSingleFile(url, download_id, &status, files)
		}
		di := DownloadsInfo[download_id]
		di.End_time = "2019-08-18 09:46:42.89588 +0530 IST m=+12.135120014"
		di.Status = status
		responseid,_ := json.Marshal(Response{Id: download_id})
		w.Header().Set("Content-Type","application/json")
		w.Write(responseid)
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

func main(){
	h := http.NewServeMux()
	h.Handle("/health", HealthHandler{})
	h.Handle("/download/", StatusHandler{})
	h.Handle("/download", DownloadHandler{})
	h.HandleFunc("/browse", func(w http.ResponseWriter, r *http.Request){
		responseBody,_ := json.Marshal(DownloadsInfo)
		w.Header().Set("Content-Type","application/json")
		w.Write(responseBody)
	})
	err := http.ListenAndServe(":8002", h)
	log.Fatal(err)
}
