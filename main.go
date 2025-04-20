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

const maxMessageLength = 1024


func handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	clientIP := strings.Split(clientAddr, ":")[0]

	fmt.Printf("[+] New connection from %s at %s\n", clientAddr, time.Now().Format(time.RFC3339))

	defer func() {
		fmt.Printf("[-] Disconnected: %s at %s\n", clientAddr, time.Now().Format(time.RFC3339))
	}()

	// Open per-client log file
	logFileName := fmt.Sprintf("%s.log", clientIP)
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("[!] Failed to open log file for %s: %v\n", clientIP, err)
		return
	}
	defer logFile.Close()

	reader := bufio.NewReader(conn)

	for {
		// 🔥 Set 30-second read timeout
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))

		message, err := reader.ReadString('\n')
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Printf("[!] %s disconnected due to inactivity\n", clientAddr)
			} else if err == io.EOF {
				fmt.Printf("[!] %s disconnected (EOF)\n", clientAddr)
			} else {
				fmt.Printf("[!] Error reading from %s: %v\n", clientAddr, err)
			}
			return
		}

		trimmed := strings.TrimSpace(message)

		// 🔐 Overflow protection
		if len(trimmed) > maxMessageLength {
			trimmed = trimmed[:maxMessageLength]
			_, _ = conn.Write([]byte("Message too long. Truncated to 1024 bytes.\n"))
		}

		// Log to file
		logLine := fmt.Sprintf("[%s] %s\n", time.Now().Format(time.RFC3339), trimmed)
		if _, err := logFile.WriteString(logLine); err != nil {
			fmt.Printf("[!] Failed to write to log for %s: %v\n", clientIP, err)
		}

		// Personality mode responses
var response string
shouldClose := false

// Command protocol
if strings.HasPrefix(trimmed, "/") {
	switch {
	case trimmed == "/time":
		response := time.Now().Format("15:04:05 MST 2006-01-02")
		_, err = conn.Write([]byte(response + "\n"))
		if err != nil {
			fmt.Printf("[!] Error writing to %s: %v\n", clientAddr, err)
			return
		}
		continue

	case trimmed == "/quit":
		_, _ = conn.Write([]byte("Goodbye!\n"))
		fmt.Printf("[i] Client %s sent /quit\n", clientAddr)
		return

	case strings.HasPrefix(trimmed, "/echo "):
		msg := strings.TrimPrefix(trimmed, "/echo ")
		_, err = conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Printf("[!] Error writing to %s: %v\n", clientAddr, err)
			return
		}
		continue

	default:
		_, _ = conn.Write([]byte("Unknown command\n"))
		continue
	}
}


switch trimmed {
case "hello":
	response = "Hi there!"
case "":
	response = "Say something..."
case "bye":
	response = "Goodbye!"
	shouldClose = true
default:
	response = trimmed // default echo
}

// Send response
_, err = conn.Write([]byte(response + "\n"))
if err != nil {
	fmt.Printf("[!] Error writing to %s: %v\n", clientAddr, err)
	return
}

if shouldClose {
	fmt.Printf("[i] Closing connection with %s after 'bye'\n", clientAddr)
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
