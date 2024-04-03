package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if r.Method != http.MethodPost {
			log.Println("middleware error: method not allowed")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		now := time.Now()
		next.ServeHTTP(w, r)
		defer func() {
			log.Printf("Request took %s", time.Since(now))
		}()
	})
}

// using for learning and debug hehe
func writeInfo(w http.ResponseWriter, r *http.Request) {
	body := fmt.Sprintf("Method: %s\r\n", r.Method)
	body += "Header ===============\r\n"
	for k, v := range r.Header {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	body += "Query parameters ===============\r\n"
	log.Println(r.URL.Path)
	for k, v := range r.URL.Query() {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
		log.Println(k, v)
	}
	if err := r.ParseForm(); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	for k, v := range r.Form {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	w.Write([]byte(body))
}
