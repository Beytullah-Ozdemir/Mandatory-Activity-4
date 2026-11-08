package main

import (
	"sync"

	proto "github.com/Beytullah-Ozdemir/Mandatory-Activity-4/grpc"
)

type State int

const (
	Released State = iota
	Wanted
	Held
)

type Node struct {
	id                 int64
	state              State
	TimestampOnRequest int64
	clock              LamportClock
	mu                 sync.Mutex
	deferred           []chan struct{}
	peers              map[string]proto.NodeServiceClient
	deferred           []chan struct{}
}
