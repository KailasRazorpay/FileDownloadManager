package main

import "time"

//Structure for holding download information
type DownloadInfo struct{
	Id string	//Download Id
	Start_time time.Time	//Download start time
	End_time time.Time	//Download end time
	Status string	//Stores whether download in pending, succesful or failed
	Download_type string	//Stores whether downloading is serial or concurrent
	Files map[string]string	//List of download urls and their local locations
}

//Structure for accepting download urls and the type of downloading
type Payload struct{
	Type string	//Download type
	Urls []string	//List of urls to download files from
}

//Structure holding response to download request
type Response struct{
	Id string	//Download Id
}

//Structure for storing error message fields
type Error struct{
	Internal_code int	//Error code
	Message string	//Error message
}

//List of all download ids and their information
var DownloadsInfo = map[string]DownloadInfo{}

