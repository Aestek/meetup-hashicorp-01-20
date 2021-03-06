package main

import (
	"flag"
	"log"
)

func main() {
	target := flag.String("target", "", "")
	rate := flag.Int("rate", 1, "Requets per second")
	listen := flag.String("listen", ":3000", "")
	theme := flag.Bool("rogue", false, "")
	flag.Parse()

	client := NewClient(*target, *rate)
	server := NewServer(*listen, client.Stats, *theme)

	go client.Run()
	err := server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
