package main

import (
	//"../Driver"
	"fmt"
	"net"
	"os"
)

func String_to_orders(Orders1 string) [4][3]int {

	//var Orders int = [12] "000000000000" //UUUUDDDDCCCC (U = orders button_up | D = orders button_down | C = orders button_command)
	var Orders_list = [4] /*N_FLOORS*/ [3] /*N_BUTTONS*/ int{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}
	//fmt.Println(Orders1)
	for i := 0; i < 4; i++ {
		if Orders1[i] == byte(49) {
			Orders_list[i][0] = 1
		} else if Orders1[i] == byte(48) {
			Orders_list[i][0] = 0
		} else {
			fmt.Println("Button_Up " + string(i) + " has an illegal value")
		}
	}
	for j := 4; j < 8; j++ {
		if Orders1[j] == byte(49) {
			Orders_list[j-4][1] = 1
		} else if Orders1[j] == byte(48) {
			Orders_list[j-4][1] = 0
		} else {
			fmt.Println("Button_Down " + string(j) + " has an illegal value")
		}
	}
	for k := 8; k < 12; k++ {
		if Orders1[k] == byte(49) {
			Orders_list[k-8][2] = 1
		} else if Orders1[k] == byte(48) {
			Orders_list[k-8][2] = 0
		} else {
			fmt.Println("Button_Command " + string(k) + "has an illegal value")
		}
	}
	return Orders_list
}

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
	var buf2 [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		//fmt.Println(string(buf[0:]))
		_, err2 := conn.Write(buf[0:n])
		if err2 != nil {
			return
		}
		var x string = string(buf[0:])
		fmt.Println(String_to_orders(x))

		m, err3 := conn.Read(buf2[0:])
		if err3 != nil {
			return
		}
		//fmt.Println(string(buf[0:]))
		_, err4 := conn.Write(buf2[0:m])
		if err4 != nil {
			return
		}
		var x2 string = string(buf2[0:])
		fmt.Println(String_to_orders(x2))
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
