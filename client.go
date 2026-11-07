package main

import (
	"log"

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
// 2. Enter critical section
func (n *Node) enterCS() { 
	n.mu.Lock()
	n.state = Wanted 
	n.TimestampOnRequest 
	n.mu.Unlock()

	log.Printf("[Node %d] event=REQUEST_SENT lamport=%d", n.id, ts) 

	var wg sync.WaitGroup,
	for _, peer := range n.peers { 
		wg.Add(1) 
		go func(peer proto.NodeServiceClient) {
			defer wg.Done()
			reply, err := peer.RequestAccess(context.Background(), 
				&proto.Request{LogicalTimestamp: ts, NodeId: n.id}, 
				grpc.WaitForReady(true))
			if err != nil {
				log.Fatalf("[Node %d] request failed: %v", n.id, err)
			} 
			n.clock.Update(reply.LogicalTimestamp)
			log.Printf("[Node %d] event=REPLY_RECEIVED from=%d", n.id, reply.NodeId, )
		}(peer)
	} 
	wg.Wait()

	n.mu.Lock()
	n.state = Held
	n.mu.Unlock()
} 
// 3. criticalSection simulates the execution of the critical section by logging entry and exit events, along with the current Lamport timestamp.
func (n *Node) criticalSection() {
	log.Printf("[Node %d] event=ENTERED_CS lamport=%d", n.id, n.clock.Get()) 
	time.Sleep(2 * time.Second)
	log.Printf("[Node %d] event=LEAVING_CS", n.id)
} 

func (n *Node) exitCS() { 
	n.muLock()()
	n.state = Released
	waiting := m.deferred 
	n.deferred = nil
	n.mu.Unlock() 

	for _, 