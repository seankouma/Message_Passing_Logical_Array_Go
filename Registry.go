package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func runRegistry() {
	messengers := make(map[uint32]net.Conn)
	go handleRequests("localhost", "8079", messengers)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ready for input")
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading")
			return
		}
		if strings.Contains(input, "connect") {
			setupRing(messengers)
		}
	}
}

func handleRequests(host string, port string, messengers map[uint32]net.Conn) {
	listen, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// close listener
	defer listen.Close()
	for {
		fmt.Println("Listening")
		conn, err := listen.Accept()
		fmt.Println("Accepted")
		fmt.Println("Node Connected")
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		fmt.Println("Size of messengers: " + strconv.Itoa(len(messengers)))
		buffer := make([]byte, 4)
		for {
			_, err := conn.Read(buffer)
			if err != nil {
				log.Fatal(err)
			}
			messageId := binary.LittleEndian.Uint32(buffer)
			messageIdString := strconv.Itoa(int(messageId))
			println("Message Id: " + messageIdString)
			switch messageId {
			case REGISTER_REQUEST:
				_, err = conn.Read(buffer)
				nodeId := binary.LittleEndian.Uint32(buffer)
				nodeIdString := strconv.Itoa(int(nodeId))
				println("Random ID: " + nodeIdString)
				_, err = conn.Read(buffer)
				port := binary.LittleEndian.Uint32(buffer)
				portString := strconv.Itoa(int(port))
				println("Port: " + portString)
				messengers[nodeId] = conn
			}
		}
	}
}
func setupRing(messengers map[uint32]net.Conn) {
	keys := make([]uint32, 0, len(messengers))
	for k := range messengers {
		keys = append(keys, k)
		fmt.Println(strconv.Itoa(len(messengers)))
	}
	for i := 0; i < len(keys)-1; i++ {
		port, _ := strconv.Atoi(messengers[keys[i+1]].LocalAddr().String()[strings.Index(messengers[keys[i+1]].LocalAddr().String(), ":"):])
		connect := GetConnectionsDirective(keys[i+1], uint32(port))
		messengers[keys[i]].Write(GetConnectionsDirectiveBytes(connect))
	}
	port, _ := strconv.Atoi(messengers[keys[0]].LocalAddr().String()[strings.Index(messengers[keys[0]].LocalAddr().String(), ":"):])
	connect := GetConnectionsDirective(keys[0], uint32(port))
	messengers[keys[len(keys)-1]].Write(GetConnectionsDirectiveBytes(connect))
}
