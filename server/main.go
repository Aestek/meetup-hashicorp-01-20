package main

import (
	"flag"
	"log"
)

func main() {
	listen := flag.String("listen", ":8001", "")
	flag.Parse()

	server := NewServer(*listen)
	err := server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
