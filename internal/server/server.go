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

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/profile", profileHandler)

	mux.HandleFunc("/process/profile", processProfileHandler)
	mux.HandleFunc("/process/register", processRegisterHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/exit", exitHandler)
	handler := authMiddleware(loggingMiddleware(headersMiddleware(mux)))
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		next.ServeHTTP(w, r)
	})
}
