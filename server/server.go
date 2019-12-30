package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/JJ/pigo"
	_ "github.com/aestek/meetup-hashicorp-01-20/server/statik"
	"github.com/rakyll/statik/fs"
	"golang.org/x/time/rate"
)

type Stats struct {
	Requests int `json:"requests,omitempty"`
}

type Server struct {
	listen       string
	currentStats Stats
}

func NewServer(listen string) *Server {
	return &Server{
		listen: listen,
	}
}

func (s *Server) Run() error {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	var limiter = rate.NewLimiter(20, 3)

	// Serve the contents over HTTP.
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(statikFS)))
	http.Handle("/", http.RedirectHandler("/public/index.html", http.StatusTemporaryRedirect))
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		s.currentStats.Requests++
		pi := pigo.Pi(50)
		w.Write([]byte("Pi is " + pi))
	})

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(s.currentStats)
	})

	log.Printf("listening on %s", s.listen)
	return http.ListenAndServe(s.listen, nil)
}
