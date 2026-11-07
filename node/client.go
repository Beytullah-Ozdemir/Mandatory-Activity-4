package main

// client.go — what a node does when it wants to ENTER the critical section:
// send a request to every peer, wait for all replies, enter CS, exit CS.
