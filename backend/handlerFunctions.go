package backend

import (
	"html/template"
	//"log"
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	//"strconv"
)

type Formdata struct {
	Input string `json:"input"`
	Color string `json:"color"`
	Font  string `json:"font"`
	File  string `json:"file"`
	Art   string `json:"art"`
}

var formData Formdata
var art string

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.URL.Path != "/" {
		errorTmpl, err := template.ParseFiles("template/404.html")
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		errorTmpl.Execute(w, nil)
		return
	}
	if r.URL.Path == "/" && r.Method == "POST" {
		http.Error(w, "400 Bad Request: This resource only supports the POST method", http.StatusBadRequest)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	//postIP = r.RemoteAddr
	defer func() {
		if err := recover(); err != nil {
			// Return the error status
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)

		}
	}()
	if r.URL.Path != "/ascii-art" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if r.URL.Path == "/ascii-art" && r.Method == "GET" {
		http.Error(w, "400 Bad Request: This resource only supports the POST method", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	input := formData.Input
	if len(input) > 20000 {
		http.Error(w, "500 Internal Server Error: The process takes exceeds the allowed time limit", http.StatusInternalServerError)
	}
	// Parse the JSON data

	// Process the data
	fmt.Printf("Received form data:\nInput: %s\nColor: %s\nFont: %s\n", formData.Input, formData.Color, formData.Font)
	art = AsciiArt(formData.Input, formData.Font)

	// Send a response
	response := map[string]string{"art": art}
	// createFile(art, "ascii-art."+formData.File)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside downoad")
	fmt.Println(formData.File)

	fileExt, contentType := getFileInfo(formData.File)

	fileName := "asciiart" + fileExt
	file, err := os.Create(fileName)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	_, err = file.WriteString(art)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Set the appropriate headers for the file download
	fileInfo, _ := file.Stat()
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileInfo.Name()))
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// Copy the file to the response writer
	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
	os.Remove(fileName)
}

func GetAscii(ch rune, i int, f string) string {
	if f != "standard" && f != "shadow" && f != "thinkertoy" {
		return ""
	}
	if i == 9 {
		return " "
	}
	chAs := (int(ch) - 32) * 9
	readFile, err := os.Open("./backend/fonts/" + f + ".txt")

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)
	counter := 0
	line := ""

	for fileScanner.Scan() {
		if counter == chAs+i {
			line += fileScanner.Text()
		}
		counter++
	}

	readFile.Close()
	return line
}

func AsciiArt(s string, f string) string {
	output := ""
	text := s
	textList := strings.Split(text, "\n")
	for _, str := range textList {
		for i := 1; i <= 8; i++ {
			for _, char := range str {
				output += GetAscii(char, i, f)
			}
			if str != "" {
				output += "\n"
			}

		}
		if str == "" {
			output += "\n"
		}
	}
	return output
}

func getFileInfo(docType string) (fileExt, contentType string) {
	switch docType {
	case "Plain Text":
		return ".txt", "text/plain"
	case "Rich Text Format":
		return ".rtf", "application/rtf"
	case "Markdown":
		return ".md", "text/markdown"
	case "Word":
		return ".doc", "application/msword"
	default:
		return ".txt", "text/plain"
	}
}
