/*package main

import (
	"net"
	"fmt"
	"bufio"
	"strings"
)

func main(){

	fmt.Println("Launching server...")
	// listen on all interfaces
	ln, _ := ln.listen("tcp",":8081")
	// accept connection on port
	conn, _ := ln.Accept()
	// run loop forever (or until ctrl-c)
	for {
		// will listen for message to process ending in newline (\n)
		message, _ := bufio.NewReader(conn).ReadString(\n)
		// output message received
		fmt.print("Message received: ", string(message))
		// sample process for string received
		newmessage := strings.ToUpper(message)
		// send new string back to client
		conn.write([]byte(newmessage + "\n"))
	}
}*/

/* SimpleEchoServer / Multithread
 */
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	service := ":1201"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		handleClient(conn)
		conn.Close() // we're finished
	}
}

func handleClient(conn net.Conn) {
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		fmt.Println(string(buf[0:]))
		_, err2 := conn.Write(buf[0:n])
		if err2 != nil {
			return
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
