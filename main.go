package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
	"flag"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	clientIP := strings.Split(clientAddr, ":")[0] // Get IP only (remove port)

	fmt.Printf("[+] New connection from %s at %s\n", clientAddr, time.Now().Format(time.RFC3339))

	defer func() {
		fmt.Printf("[-] Disconnected: %s at %s\n", clientAddr, time.Now().Format(time.RFC3339))
	}()

	// Open log file for this client
	logFileName := fmt.Sprintf("%s.log", clientIP)
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("[!] Failed to open log file for %s: %v\n", clientIP, err)
		return
	}
	defer logFile.Close()

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("[!] %s disconnected (EOF)\n", clientAddr)
			} else {
				fmt.Printf("[!] Error reading from %s: %v\n", clientAddr, err)
			}
			return
		}

		trimmed := strings.TrimSpace(message)

		// Log to file
		logLine := fmt.Sprintf("[%s] %s\n", time.Now().Format(time.RFC3339), trimmed)
		if _, err := logFile.WriteString(logLine); err != nil {
			fmt.Printf("[!] Failed to write to log for %s: %v\n", clientIP, err)
		}

		// Echo back
		_, err = conn.Write([]byte(trimmed + "\n"))
		if err != nil {
			fmt.Printf("[!] Error writing to %s: %v\n", clientAddr, err)
			return
		}
	}
}

func main() {
	// Define a port flag with default value 4000
	port := flag.String("port", "4000", "TCP port to listen on")
	flag.Parse()

	address := fmt.Sprintf(":%s", *port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("[!] Failed to start server on port %s: %v\n", *port, err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("[✓] Server listening on port %s\n", *port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[!] Error accepting:", err)
			continue
		}
		go handleConnection(conn) // concurrency already added!
	}
}
