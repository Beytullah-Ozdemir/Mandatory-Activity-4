package main

import(
	"context"
	"log"
	"net"
	"fmt"
	proto "github.com/Beytullah-Ozdemir/Mandatory-Activity-4/grpc"
	"google.golang.org/grpc"
)
// server.go — what a node does when it RECEIVES a request from another node.
// (Ricart–Agrawala: reply now, or defer the reply until we leave the CS.)

type NodeServer struct{
proto.UnimplementedNodeServiceServer
node *Node
}



func(s *NodeServer) start_server(port string) {

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Server did not start: %v", err)
		//%v indsætter den konkrete fejl og fejlen bliver logget
	}

	proto.RegisterNodeServiceServer(grpcServer, s)
	log.Printf("[Server] event=STARTUP port=%s", port) 
	
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}

func(s *NodeServer) RequestAccess(ctx context.Context, req *proto.Request)(*proto.Reply, error){

	s.node.mu.Lock()
	s.node.clock.Update(req.GetLogicalTimestamp())
	log.Printf("[Node %d] event=REQUEST_RECEIVED from=%d lamport=%d request_time=%d", s.node.id, req.GetNodeId(), s.node.clock, req.GetLogicalTimestamp())
	switch s.node.state {
	
	case Released:
		s.node.mu.Unlock()
		 return &proto.Reply{}, nil 
	
	case Wanted: 
		
		if req.GetLogicalTimestamp() < s.node.TimestampOnRequest{
			s.node.mu.Unlock() 
			return &proto.Reply{}, nil 

		}else if req.GetLogicalTimestamp() == s.node.TimestampOnRequest && req.GetNodeId() < s.node.id{
			 s.node.mu.Unlock()
			return &proto.Reply{}, nil 

		}else{
			// udsæt svar
		}

		log.Println("Node is WANTED")
	
	case Held: 
		log.Println("Node is HELD")
	}
	s.node.mu.Unlock()


	return &proto.Reply{req}, nil

}


func (n *Node) makeReply(req *proto.Request) *proto.Reply {
	log.Printf ("[Node %d] event=REPLY_SENT to=%d lamport=%d", s.node.id, req.GetNodeId(), s.node.clock.Now())
	return &proto.Reply{}
}

	
