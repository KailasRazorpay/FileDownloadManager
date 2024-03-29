package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

//Maximum number of concurrent goroutines allowed
var numberOfLinks int = 2

//Handler for responding to /health
type HealthHandler struct{}

func (re HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(200)
	fmt.Fprintln(w,"OK")
}

//Handler for responding to /downloads/<download-id>
type StatusHandler struct{}

func (st StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	segments := strings.Split(r.URL.Path, "/")
	id := segments[2]

	//Checks if download id is valid
	if _ , val := DownloadsInfo[id]; val{

		//Prints download information to console
		fmt.Println(DownloadsInfo[id])

		//Converts download information to json result
		responseBody,err := json.Marshal(DownloadsInfo[id])
		if err != nil{
			panic(err)
		}

		//Returns download information in response
		w.Header().Set("Content-Type","application/json")
		w.Write(responseBody)
	} else{	//If download received is invalid

		//Prepares error message
		errorBody := Error{
			Internal_code : 4001,
			Message : "unknown download id",
		}

		//Returns error message
		w.Header().Set("Content-Type","application/json")
		error,_ := json.Marshal(errorBody)
		w.Write(error)
	}
}

//Handler for responding to download requests at /downloads
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
		serialDownload(w, r, payload)

	} else if(payload.Type == "concurrent"){
		concurrentDownload(w, r, payload)

	} else {

		//Prepares error message
		errorBody := Error{
			Internal_code : 4002,
			Message : "unknown type of download",
		}

		//Returns error message
		w.Header().Set("Content-Type","application/json")
		error,_ := json.Marshal(errorBody)
		w.Write(error)
	}
}

