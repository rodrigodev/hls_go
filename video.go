package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/xfrr/goffmpeg/transcoder"
)

const mediaRoot = "static/media"
const outputFile = "test.m3u8"
const uploadDir = "upload/"
const maxMemory = 32 << 20

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	mId, err := strconv.Atoi(vars["mId"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	segName, ok := vars["segName"]
	if !ok {
		mediaBase := getMediaBase(mId)
		m3u8Name := "index.m3u8"
		serveHlsM3u8(w, r, mediaBase, m3u8Name)
	} else {
		mediaBase := getMediaBase(mId)
		serveHlsTs(w, r, mediaBase, segName)
	}
}

func getMediaBase(mId int) string {
	return fmt.Sprintf("%s/%d", mediaRoot, mId)
}

func serveHlsM3u8(w http.ResponseWriter, r *http.Request, mediaBase, m3u8Name string) {
	log.Print("m3u8 hit")

	mediaFile := fmt.Sprintf("%s/hls/%s", mediaBase, m3u8Name)
	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "application/x-mpegURL")
}

func serveHlsTs(w http.ResponseWriter, r *http.Request, mediaBase, segName string) {
	log.Printf("segment hit %s", segName)

	mediaFile := fmt.Sprintf("%s/hls/%s", mediaBase, segName)
	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "video/MP2T")
}

func mainUpload(w http.ResponseWriter, r *http.Request)  {
	http.ServeFile(w, r, "static/upload.html")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(maxMemory)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile(uploadDir + handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	go func(path string) {
		transcodeToHls(path)
	}(handler.Filename)

	return
}

func transcodeToHls(p string) {
	trans := new(transcoder.Transcoder)

	err := trans.Initialize(p, uploadDir + outputFile)
	if err != nil {
		log.Fatal(err)
	}

	done := trans.Run(true)
	progress := trans.Output()

	for msg := range progress {
		fmt.Println(msg)
	}

	err = <-done
	if err != nil {
		log.Fatal(err)
	}
}