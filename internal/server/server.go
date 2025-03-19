package server

import (
	"log"
	"os"

	"net/http"

	_ "github.com/lib/pq"
)

func Start() {
	var port string = os.Getenv("APP_PORT")

	mux := http.NewServeMux()

	mux.HandleFunc("/", authMiddleware(homeHandler))
	mux.HandleFunc("/form", formHandler)
	mux.HandleFunc("/process", processHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/exit", exitHandler)

	handler := loggingMiddleware(headersMiddleware(mux))
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Println("starting server..")
	if err := http.ListenAndServe("0.0.0.0:"+port, handler); err != nil {
		log.Fatal(err)
	}
}

// Middleware для логирования запросов
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Запрос: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Middleware для установки header
func headersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET POST PUT OPTIONS CONNECT HEAD")
		next.ServeHTTP(w, r)
	})
}
