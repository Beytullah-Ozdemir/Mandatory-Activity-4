package main

import (
	"context"
	"log"
	"net"

	proto "github.com/Beytullah-Ozdemir/Mandatory-Activity-4/grpc"
	"google.golang.org/grpc"
)

type NodeServer struct {
	proto.UnimplementedNodeServiceServer
	node *Node
}

func (s *NodeServer) start_server(port string) {

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

func (s *NodeServer) RequestAccess(ctx context.Context, req *proto.Request) (*proto.Reply, error) {

	s.node.mu.Lock()
	s.node.clock.Update(req.GetLogicalTimestamp())
	log.Printf("[Node %d] event=REQUEST_RECEIVED from=%d lamport=%d request_time=%d", s.node.id, req.GetNodeId(), s.node.clock.Now(), req.GetLogicalTimestamp())
	switch s.node.state {

	case Released:
		s.node.mu.Unlock()
		return &proto.Reply{}, nil

	case Wanted:

		if req.GetLogicalTimestamp() < s.node.TimestampOnRequest {
			s.node.mu.Unlock()
			return &proto.Reply{}, nil

		} else if req.GetLogicalTimestamp() == s.node.TimestampOnRequest && req.GetNodeId() < s.node.id {
			s.node.mu.Unlock()
			return &proto.Reply{}, nil

		} else {
			return s.waitToReply()
		}

	case Held:
		return s.waitToReply()
	}
	return &proto.Reply{}, nil

}

func (n *Node) makeReply(req *proto.Request) *proto.Reply {
	n.clock.Tick()
	log.Printf("[Node %d] event=REPLY_SENT to=%d lamport=%d", n.id, req.GetNodeId(), n.clock.Now())
	return &proto.Reply{}
}

func (s *NodeServer) waitToReply() (*proto.Reply, error) {
	ready := make(chan struct{})
	s.node.deferred = append(s.node.deferred, ready)
	s.node.mu.Unlock()
	<-ready
	return &proto.Reply{}, nil
}
