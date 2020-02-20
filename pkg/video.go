package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/xfrr/goffmpeg/transcoder"
)

const mediaPath = "media"
const uploadPath = "upload"
const indexFileName = "index.m3u8"
const maxMemory = 32 << 20

func getMediaBase(mId int) string {
	return fmt.Sprintf("%s/%d", mediaPath, mId)
}

func videoEmbedHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func uploadFormHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/upload.html")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		fmt.Println(err)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Fprintf(w, "%v", handler.Header)

	uploadPath := fmt.Sprintf("%s/%s", uploadPath, handler.Filename)

	f, err := os.OpenFile(uploadPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Println(err)
		return
	}

	transcodeToHls(handler.Filename)
}

func transcodeToHls(filePath string) {
	trans := new(transcoder.Transcoder)

	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}

	m3u8Path := fmt.Sprintf("%s/%s", mediaPath, id)
	m3u8File := fmt.Sprintf("%s/%s", m3u8Path, indexFileName)
	err = os.MkdirAll(m3u8Path, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(m3u8File)

	err = trans.Initialize(filePath, m3u8File)
	if err != nil {
		log.Print("Fail init")
		log.Fatal(err)
	}

	done := trans.Run(true)
	progress := trans.Output()

	for msg := range progress {
		fmt.Println(msg)
	}

	err = <-done
	if err != nil {
		log.Print("Fail Processing")
		log.Fatal(err)
	}
}
