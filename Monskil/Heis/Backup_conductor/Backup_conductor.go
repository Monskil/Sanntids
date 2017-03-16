package Backup_conductor

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"time"
)

func Backup_conductor() {
	for {
		primal_allive := true

		addr, _ := net.ResolveUDPAddr("udp", "localhost"+":40016")

		conn_secondary, err := net.ListenUDP("udp", addr)

		if err != nil {
			log.Fatal("No local UDP adress", err)
		}
		connection_ch := make(chan bool)
		go listen_UDP_primal(conn_secondary, connection_ch)

		last_connection_time := time.Now()
		for primal_allive {
			select {
			case <-connection_ch:
				last_connection_time = time.Now()
			default:
				if time.Since(last_connection_time).Seconds() > 2.0 {
					primal_allive = false
					conn_secondary.Close()
				}
			}

		}
		terminal := exec.Command("gnome-terminal", "-x", "go", "run", "Main-2.go")
		terminal.Output()

		conn_primal, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			log.Fatal("Not able to DialUDP")
		}

		defer conn_primal.Close()
		time.Sleep(1 * time.Second)
		for {
			msg := []byte("Primal alive")
			_, err = conn_primal.Write(msg)
			if err != nil {
				fmt.Println("Something went wrong!")
			}

			time.Sleep(1 * time.Second)
		}
	}

}

func listen_UDP_primal(conn *net.UDPConn, connection_ch chan bool) {
	for {

		msg := make([]byte, 1024)
		_, _, err := conn.ReadFromUDP(msg)

		//fmt.Println(a)
		if err != nil {
			fmt.Println("Something went wrong!")
		}
		connection_ch <- true
	}
}
