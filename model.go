package main

import "time"

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

