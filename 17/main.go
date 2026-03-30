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

func main() {
	timeoutSeconds := flag.Int("timeout", 10, "timeout in seconds")
	flag.Parse()

	timeout := time.Second * time.Duration(*timeoutSeconds)

	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("usage: telnet [-t <timeout in seconds>] <host> <port>")
		os.Exit(1)
	}

	host := args[0]
	port := args[1]

	addr := net.JoinHostPort(host, port)

	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		fmt.Println("failed to connect:", err)
		os.Exit(1)
	}
	defer func() {
		_ = conn.Close()
	}()

	fmt.Println("Connected to", addr)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			_, err := fmt.Fprintf(conn, "%s\r\n", scanner.Text())
			if err != nil {
				fmt.Println("failed to send data:", err)
				return
			}
		}
	}()

	if _, err := io.Copy(os.Stdout, conn); err != nil {
		fmt.Println("failed to read:", err)
	}

}
