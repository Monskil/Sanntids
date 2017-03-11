package Network

import (
	//"../Driver"
	"bufio"
	"fmt"
	"net"
	"os"
)

///IP_ADRESSE 129.241.187.tall

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

var allClients map[*Client]int

type Client struct {
	// incoming chan string
	outgoing   chan string
	reader     *bufio.Reader
	writer     *bufio.Writer
	conn       net.Conn
	connection *Client
}

func (client *Client) Read() {
	for {
		line, err := client.reader.ReadString('\n')
		if err == nil {
			if client.connection != nil {
				client.connection.outgoing <- line
			}
			fmt.Println(line)
		} else {
			break
		}
	}
	client.conn.Close()
	delete(allClients, client)
	if client.connection != nil {
		client.connection.connection = nil
	}
	client = nil
}

func (client *Client) Write() {
	for data := range client.outgoing {
		client.writer.WriteString(data)
		client.writer.Flush()
	}
}

func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}

func NewClient(connection net.Conn) *Client {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	client := &Client{
		// incoming: make(chan string),
		outgoing: make(chan string),
		conn:     connection,
		reader:   reader,
		writer:   writer,
	}
	client.Listen()

	return client
}

func Network_server_main( /*New_order bool*/ ) {
	allClients = make(map[*Client]int)
	listener, _ := net.Listen("tcp", ":1201")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}

		client := NewClient(conn)
		for clientList, _ := range allClients {
			if clientList.connection == nil {
				client.connection = clientList
				clientList.connection = client
				fmt.Println("Connected")
			}
			allClients[client] = 1
			fmt.Println(len(allClients))
		}
		go handleClient(conn /*, New_order*/)
	}
}

//var lol2 bool = Driver.Bursdagskvinn()
var Server_list = [4] /*N_FLOORS*/ [3] /*N_BUTTONS*/ int{
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
}

func handleClient(conn net.Conn /*, New_order bool*/) {
	for {
		defer conn.Close()
		var buf [12] /*512*/ byte
		for {
			n, err := conn.Read(buf[0:])
			if err != nil {
				return
			}
			_, err2 := conn.Write(buf[0:n])
			if err2 != nil {
				return
			}
			var x string = string(buf[0:]) // + string('\n')
			fmt.Println(x)
			Server_list = /*fmt.Println(*/ String_to_orders(x) //)
			fmt.Println(String_to_orders(x))
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
