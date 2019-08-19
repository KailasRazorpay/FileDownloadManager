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
	"time"
)

type DownloadInfo struct{
	Id string
	Start_time time.Time
	End_time time.Time
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

type Error struct{
	Internal_code int
	Message string
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
	if _ , val := DownloadsInfo[id]; val{
		fmt.Println(DownloadsInfo[id])
		responseBody,err := json.Marshal(DownloadsInfo[id])
		if err != nil{
			panic(err)
		}
		w.Header().Set("Content-Type","application/json")
		w.Write(responseBody)
	} else{
		errorBody := Error{
			Internal_code : 4001,
			Message : "unknown download id",
		}
		w.Header().Set("Content-Type","application/json")
		error,_ := json.Marshal(errorBody)
		w.Write(error)
	}
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
		status := "PENDING"
		download_id := xid.New().String()
		files := make(map[string]string)
		start_time := time.Now()
		for _,url := range payload.Urls{
			downloadSingleFile(url, download_id, &status, files)
		}
		if status == "PENDING" {
			status = "SUCCESSFUL"
		}
		end_time := time.Now()
		DownloadsInfo[download_id] = DownloadInfo{
			Id : download_id,
			Start_time : start_time,
			End_time : end_time,
			Status : status,
			Download_type : payload.Type,
			Files : files,
		}
		responseid,_ := json.Marshal(Response{Id: download_id})
		w.Header().Set("Content-Type","application/json")
		w.Write(responseid)
	} else if(payload.Type == "concurrent"){
		status := "PENDING"
		download_id := xid.New().String()
		responseid,_ := json.Marshal(Response{Id: download_id})
		w.Header().Set("Content-Type","application/json")
		w.Write(responseid)
		files := make(map[string]string)
		start_time := time.Now()
		for _,url := range payload.Urls{
			go downloadSingleFile(url, download_id, &status, files)
		}
		if status == "PENDING" {
			status = "SUCCESSFUL"
		}
		end_time := time.Now()
		DownloadsInfo[download_id] = DownloadInfo{
			Id : download_id,
			Start_time : start_time,
			End_time : end_time,
			Status : status,
			Download_type : payload.Type,
			Files : files,
		}
	} else {
		errorBody := Error{
		Internal_code : 4002,
		Message : "unknown type of download",
		}
		w.Header().Set("Content-Type","application/json")
		error,_ := json.Marshal(errorBody)
		w.Write(error)
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
	h.Handle("/downloads/", StatusHandler{})
	h.Handle("/downloads", DownloadHandler{})
	h.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request){
		responseBody,_ := json.Marshal(DownloadsInfo)
		w.Header().Set("Content-Type","application/json")
		w.Write(responseBody)
	})
	err := http.ListenAndServe(":8002", h)
	log.Fatal(err)
}
