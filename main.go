package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("[+] New connection from %s at %s\n", clientAddr, time.Now().Format(time.RFC3339))

	defer func() {
		fmt.Printf("[-] Disconnected: %s at %s\n", clientAddr, time.Now().Format(time.RFC3339))
	}()

	reader := bufio.NewReader(conn)

	for {
		// Read until newline
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("[!] %s disconnected (EOF)\n", clientAddr)
			} else {
				fmt.Printf("[!] Error reading from %s: %v\n", clientAddr, err)
			}
			return
		}

		// Trim input
		trimmed := strings.TrimSpace(message)

		// Echo it back
		_, err = conn.Write([]byte(trimmed + "\n"))
		if err != nil {
			fmt.Printf("[!] Error writing to %s: %v\n", clientAddr, err)
			return
		}
	}
}



func main(){

	listener, err:=net.Listen("tcp", ":4000")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Server listening on :4000")

	for{
		conn, err := listener.Accept()
		if err != nil{
			fmt.Println("Error accepting:", err)
			continue
		}
		go handleConnection(conn)
	}
	

}

