package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
