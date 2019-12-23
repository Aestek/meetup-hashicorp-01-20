package main

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/aestek/meetup-hashicorp-01-20/client/statik"
	"github.com/rakyll/statik/fs"
)

type Server struct {
	listen       string
	stats        chan Stats
	currentStats Stats
}

func NewServer(listen string, stats chan Stats) *Server {
	return &Server{
		listen: listen,
		stats:  stats,
	}
}

func (s *Server) Run() error {

	go func() {
		for stats := range s.stats {
			s.currentStats = stats
		}
	}()

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	// Serve the contents over HTTP.
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(statikFS)))
	http.Handle("/", http.RedirectHandler("/public/index.html", http.StatusTemporaryRedirect))

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(s.currentStats)
	})

	return http.ListenAndServe(s.listen, nil)
}
