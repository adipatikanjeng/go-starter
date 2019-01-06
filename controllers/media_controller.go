package controllers

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"rest-api/utils/caching"
	"strconv"
	"time"
)

type MediaController struct {
	DB    *sql.DB
	Cache caching.Cache
}

func NewMediaController(db *sql.DB, c caching.Cache) *MediaController {
	return &MediaController{
		DB:    db,
		Cache: c,
	}
}
func (mc *MediaController) Upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./uploads/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}
