package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func routes(){

	//Creates handler to handle all the different routes
	h := http.NewServeMux()

	//Returns status 200 if get request is made on /health url
	h.Handle("/health", HealthHandler{})

	//Returns download information corresponding to a particular id or else returns an error message
	h.Handle("/downloads/", StatusHandler{})

	//Downloads files from urls included in post request to /downloads url and creates an entry in Download information
	//Returns error if mrthod is not post
	h.Handle("/downloads", DownloadHandler{})

	//Returns information about all downloads
	h.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request){
		responseBody,_ := json.Marshal(DownloadsInfo)
		w.Header().Set("Content-Type","application/json")
		w.Write(responseBody)
	})

	//Starts http server on port 8002
	err := http.ListenAndServe(":8002", h)
	log.Fatal(err)
}