package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func flagParse() (*int, string) {
	timeout := flag.Int("timeout", 10, "connection timeout in seconds")
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--timeout=sec] host port\n", os.Args[0])
		os.Exit(1)
	}
	host := flag.Arg(0)
	port := flag.Arg(1)

	address := net.JoinHostPort(host, port)

	return timeout, address
}

func main() {
	timeout, address := flagParse()

	conn, err := net.DialTimeout("tcp", address, time.Duration(*timeout)*time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect %s:%v", address, err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Printf("Connected to %s\n", address)

	done := make(chan struct{})

	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Connection closed: %v", err)
		}
		close(done)
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text() + "\n"
			_, err := conn.Write([]byte(text))
			if err != nil {
				break
			}
		}
		conn.Close()
	}()

	<-done
	fmt.Fprintln(os.Stderr, "Disconnected")
}
