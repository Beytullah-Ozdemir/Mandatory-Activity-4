package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {

	//Flag
	id := flag.Int("id", 1, "unique id for this node")
	port := flag.String("port", "5001", " The port this node listen on")
	peersFlag := flag.String("peers", "", "ports on the other nodes") // fx 5002,5003
	flag.Parse()

	// Log terminal and file
	logFile, err := s.OpenFile(fmt.Sprintf("node-%d.log", *id), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("could not open (THE LOG FILE) %v", err)

	}

	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}
