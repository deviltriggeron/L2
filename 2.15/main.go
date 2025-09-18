package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func execute(command []string) {
	pipeline := parsePipeline(command)
	for _, p := range pipeline {
		switch p[0] {
		case "cd":
			dir := os.Getenv("HOME")
			if len(p) > 1 {
				dir = p[1]
			}
			if err := os.Chdir(dir); err != nil {
				fmt.Fprintln(os.Stderr, "cd:", err)
			}
		case "pwd":
			dir, err := os.Getwd()
			if err != nil {
				fmt.Fprintln(os.Stderr, "pwd:", err)
			} else {
				fmt.Println(dir)
			}
		case "echo":
			fmt.Println(strings.Join(p[1:], " "))
		case "kill":
			if len(p[1]) < 2 {
				fmt.Fprintln(os.Stderr, "kill: pid required")
				return
			}
			pid, _ := strconv.Atoi(p[1])
			proc, err := os.FindProcess(pid)
			if err != nil {
				fmt.Fprintln(os.Stderr, "kill:", err)
				return
			}
			proc.Signal(syscall.SIGTERM)
		case "ps":
			cmd := exec.Command("ps")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
		case "ls":
			cmd := exec.Command("ls")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
		}
	}

}

func parsePipeline(tokens []string) [][]string {
	var pipeline [][]string
	var cur []string

	for _, tok := range tokens {
		if tok == "|" {
			pipeline = append(pipeline, cur)
			cur = []string{}
		} else {
			cur = append(cur, tok)
		}
	}

	if len(cur) > 0 {
		pipeline = append(pipeline, cur)
	}

	return pipeline
}

func tokenize(line string) []string {
	return strings.Fields(line)
}

func printPrompt() {
	wd, _ := os.Getwd()
	fmt.Printf("%s$ ", wd)
}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)

	go func() {
		for range sigCh {
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		printPrompt()
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			fmt.Println("error:", err)
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		tokens := tokenize(line)

		execute(tokens)

		time.Sleep(1 * time.Millisecond)
	}
}
