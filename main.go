package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {

	//Flag
	id := flag.Int("id", 1, "unique id for this node")
	port := flag.String("port", "5001", " The port this node listen on")
	peersFlag := flag.String("peers", "", "ports on the other nodes") // fx 5002,5003
	flag.Parse()

	// Log terminal and file
	logFile, err := os.OpenFile(fmt.Sprintf("node-%d.log", *id), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("could not open (THE LOG FILE) %v", err)

	}

	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	node := &Node{id: int64(*id)}

	// start gRPC-Serveren
	go (&NodeServer{node: node}).start_server(*port)

	var peerPorts []string
	for _, p := range strings.Split(*peersFlag, ",") {
		if p = strings.TrimSpace(p); p != "" {
			peerPorts = append(peerPorts, p)
		}

	}
	node.connectToPeers(peerPorts)

	go node.run(3) // request the CS 3 times, then just keep answering

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Printf("[Node %d] event=SHUTDOWN lamport=%d", node.id, node.clock.Now())

}
