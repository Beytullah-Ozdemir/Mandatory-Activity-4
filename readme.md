# Mandatory-Activity-4
Mandatory Activity 4 - Distributed Systems, IT University of Copenhagen

A small distributed system in Go. Several nodes share one critical section, and only one
node can be inside it at a time. The nodes follow the Ricart-Agrawala algorithm and talk to
each other with gRPC (every node is both a server and a client).


Requirements:
* Go 1.27.0 (specified in go.mod)
* Dependencies are managed through go.mod and go.sum
* protoc is only needed if grpc/node.proto is changed (the generated files are in the repository)

The generated files grpc/node.pb.go and grpc/node_grpc.pb.go are already in the repository, so you do not need protoc to run the system.

## Project Structure
```text
Mandatory-Activity-4/
├── grpc/
│   ├── node.proto
│   ├── node.pb.go
│   └── node_grpc.pb.go
├── .gitignore
├── client.go
├── go.mod
├── go.sum
├── lamport.go
├── main.go
├── node.go
├── node-1.log
├── node-2.log
├── node-3.log
├── readme.md
└── server.go
```


## How to Run

# Terminal 1
go run . -id 1 -port 5001 -peers 5002,5003

# Terminal 2
go run . -id 2 -port 5002 -peers 5001,5003

# Terminal 3
go run . -id 3 -port 5003 -peers 5001,5002

Each node asks for the critical section a few times at random intervals. When a node has finished, it logs event=DONE and keeps answering the others.
Press Ctrl+C in every terminal once all three nodes have logged DONE.


## How It Works

The system uses the Ricart-Agrawala algorithm to coordinate access to the critical section.
Each node can be in one of three states:

* RELEASED: The node is not requesting access to the critical section.

* WANTED: The node is waiting for permission from the other nodes.

* HELD: The node is currently inside the critical section.

When a node wants to enter the critical section, it sends a request to all other nodes containing its Lamport timestamp.
A node receiving a request replies immediately if it is not interested in entering the critical section. Otherwise, it compares the requests using their Lamport timestamps and node IDs. The request with the lower timestamp has priority. If the other node has priority, the reply is sent immediately. Otherwise, the reply is deferred.
When a node leaves the critical section, it sends replies to all nodes whose requests were deferred.
A node can enter the critical section once it has received replies from all other nodes.


## Nodes logging

Each node writes log events to the terminal and to its own log file:

* node-1.log

* node-2.log

* node-3.log

The logs contain timestamps, node IDs, Lamport clock values and event information.
The log files are appended to on every run. Delete the old log files before starting a new run if you want clean logs.
The logs can be used to verify mutual exclusion. No node should enter the critical section while another node is already inside it.
