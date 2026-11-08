package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	proto "github.com/Beytullah-Ozdemir/Mandatory-Activity-4/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 1. connectToPeers establishes gRPC connections to the specified peer ports and stores the clients in the Node's peers map.
func (n *Node) connectToPeers(peerPorts []string) {
	n.peers = make(map[string]proto.NodeServiceClient)
	for _, port := range peerPorts {
		conn, err := grpc.NewClient("localhost:"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("[Node %d] could not connect to %s: %v", n.id, port, err)
		}
		n.peers[port] = proto.NewNodeServiceClient(conn)
	}
}

// 2. enterCS asks every peer for permission and waits until ALL of them have replied.
func (n *Node) enterCS() {
	n.mu.Lock()
	n.state = Wanted
	n.TimestampOnRequest = n.clock.Tick()
	ts := n.TimestampOnRequest
	n.mu.Unlock()

	log.Printf("[Node %d] event=REQUEST_SENT lamport=%d", n.id, ts)

	var wg sync.WaitGroup
	for _, peer := range n.peers {
		wg.Add(1)
		go func(peer proto.NodeServiceClient) {
			defer wg.Done()
			// Blocks until the other node allows us in (its server.go may defer the reply).
			reply, err := peer.RequestAccess(context.Background(),
				&proto.Request{LogicalTimestamp: ts, NodeId: n.id},
				grpc.WaitForReady(true))
			if err != nil {
				log.Fatalf("[Node %d] request failed: %v", n.id, err)
			}
			lamport := n.clock.Update(reply.GetLogicalTimestamp())
			log.Printf("[Node %d] event=REPLY_RECEIVED from=%d lamport=%d", n.id, reply.GetNodeId(), lamport)
		}(peer)
	}
	wg.Wait() // wait until every peer has replied

	n.mu.Lock()
	n.state = Held
	n.mu.Unlock()
}

// 3. criticalSection simulates the execution of the critical section by logging entry and exit events, along with the current Lamport timestamp.
func (n *Node) criticalSection() {
	log.Printf("[Node %d] event=ENTERED_CS lamport=%d", n.id, n.clock.Now())
	time.Sleep(2 * time.Second)
	log.Printf("[Node %d] event=LEAVING_CS lamport=%d", n.id, n.clock.Now())
}

// 4. exitCS leaves the critical section and answers every request we deferred while we were Wanted/Held.
func (n *Node) exitCS() {
	n.mu.Lock()
	n.state = Released
	waiting := n.deferred
	n.deferred = nil
	n.mu.Unlock()

	for _, ch := range waiting {
		close(ch) // wakes up a waiting RequestAccess in server.go, which then sends its reply
	}
	log.Printf("[Node %d] event=EXITED_CS released_deferred=%d", n.id, len(waiting))
}

// 5. run keeps requesting the critical section at random intervals, forever.
func (n *Node) run() {
	for {
		time.Sleep(time.Duration(1+rand.Intn(3)) * time.Second)
		n.enterCS()
		n.criticalSection()
		n.exitCS()
	}
}
