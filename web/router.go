package web

import (
	"encoding/json"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"umutsevdi/pengui/sys"
)

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("static/index.html")
	log.Println("requested:/")
	if err != nil {
		log.Println(r.URL.String(), "#NotFound")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()
	contentType := mime.TypeByExtension(".html")
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(200)
	io.Copy(w, file)
}

func ServeStatic(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[1:]
	log.Println("requested:" + filePath)
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(r.URL.String(), "#NotFound")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()
	contentType := mime.TypeByExtension(filePath)
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(200)
	io.Copy(w, file)
}

func ServeResource(w http.ResponseWriter, r *http.Request) {
	log.Println("requested:" + r.URL.String())
	data, err := json.Marshal(sys.GetHostInfo())
	if err != nil {
		log.Println("NotAvailable")
		w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}
	contentType := mime.TypeByExtension(".json")
	w.Header().Add("Content-Type", contentType)
	w.WriteHeader(200)
	w.Write(data)
}
