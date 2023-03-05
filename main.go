package main

import (
	"net/http"
	"umutsevdi/pengui/web"
)

func main() {
	http.HandleFunc("/", web.ServeIndex)
	http.HandleFunc("/static/", web.ServeStatic)
	http.HandleFunc("/ws/init", web.ServeResource)
	http.HandleFunc("/ws/", web.ServeWs)
	http.HandleFunc("/ws/term", web.ServeTerm)
	http.ListenAndServe(":8081", nil)
}
