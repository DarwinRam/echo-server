package main

import(
	"fmt"
	"net"
	"time"
)


func handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()             // Get client address (IP:Port)
	fmt.Printf("[+] New connection from %s at %s\n",     // Log connection
		clientAddr, time.Now().Format(time.RFC3339))      // Use RFC3339 for readable timestamp

	defer func() {
		fmt.Printf("[-] Disconnected: %s at %s\n",
			clientAddr, time.Now().Format(time.RFC3339)) // Log disconnection
	}()


	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from client:", err)
			return
		}

		_, err = conn.Write(buf[:n])
	if err != nil {
		fmt.Println("Error writing to client:", err)
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

