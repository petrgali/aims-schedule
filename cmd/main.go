package main

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	Tmpl *template.Template
}

func main() {
	handler := &Handler{}
	handler.Tmpl = template.Must(template.ParseGlob("../static/*"))

	fs := http.FileServer(http.Dir("../static/"))
	mux := mux.NewRouter()
	mux.Handle("/static/", http.StripPrefix("/static", fs))
	mux.HandleFunc("/", handler.indexHandler).Methods("GET")

	srv := &http.Server{
		Addr:           ":5050",
		Handler:        logger(mux),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("server now listening on port %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("server failed to start. Reason: %s", err.Error())
	}
}

func (h *Handler) indexHandler(w http.ResponseWriter, r *http.Request) {
	h.Tmpl.ExecuteTemplate(w, "template.html", nil)
}
func logger(next http.Handler) http.Handler {
	return (http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("requested URL: '%s', method '%s'", r.URL.Path, r.Method)
		next.ServeHTTP(w, r)
	}))
}
