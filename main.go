package main

import (
	//"html/template"
	"log"
	"net/http"
	ascii "asciiart/backend"
)

func main() {
	fs := http.FileServer(http.Dir("template"))
	http.Handle("/template/", http.StripPrefix("/template/", fs))
	http.HandleFunc("/", ascii.HomeHandler)
	http.HandleFunc("/ascii-art", ascii.SubmitHandler)
    http.HandleFunc("/download", ascii.DownloadHandler)

	log.Println("Server started on http://localhost:8070")
	log.Fatal(http.ListenAndServe(":8070", nil))
}

