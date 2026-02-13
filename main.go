package main

import (
	"fileuploader/controllers"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/upload", controllers.UploadHandler)
	http.HandleFunc("/download", controllers.DownloadHandler)
	http.Handle("/", http.FileServer(http.Dir("./web")))
	fmt.Println("Server is listed")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
