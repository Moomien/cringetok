package controllers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

var allowedMood = map[string]bool{"sad": true, "happy": true, "playful": true}
var allowedExts = map[string]bool{".jpg": true, ".png": true}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	const maxfilesize = 10 << 20
	if r.ContentLength > maxfilesize {
		http.Error(w, "File is too big", http.StatusRequestEntityTooLarge)
		return
	}

	if err := r.ParseMultipartForm(maxfilesize); err != nil {
		http.Error(w, "Failed parse multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File err:"+err.Error(), http.StatusBadRequest)
		return
	}

	mood := r.FormValue("mood")
	if !allowedMood[mood] {
		http.Error(w, "Error mood", http.StatusBadRequest)
		return
	}

	ext := filepath.Ext(fileHeader.Filename)
	if !allowedExts[ext] {
		http.Error(w, "Only images are allowed", http.StatusBadRequest)
		return
	}

	defer file.Close()
	uploadDir := "./files"
	uuid := uuid.New().String()
	destPath := filepath.Join(uploadDir, uuid+filepath.Base(fileHeader.Filename))
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Failed create path: "+err.Error(), http.StatusInternalServerError)
		return
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "Failed to create file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(destFile, file); err != nil {
		http.Error(w, "Failed save file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	destFile.Close()

	videoPath := filepath.Base(GenerateVideo(destPath, mood, uuid))
	http.Redirect(w, r, "/download?file="+videoPath, http.StatusSeeOther)
}
