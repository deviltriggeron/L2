package main

import (
	"log"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("udp", host)
	if err != nil {
		log.Fatal("failed to connect:", err)
	}
	defer conn.Close()
	if err := conn.SetDeadline(
		time.Now().Add(15 * time.Second)); err != nil {
		log.Fatal("failed to set deadline: ", err)
	}
}
