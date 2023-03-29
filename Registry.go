package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
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
		} else if strings.Contains(input, "begin") {
			commence(messengers, 100000)
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
		_, err = conn.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		messageId := binary.LittleEndian.Uint32(buffer)
		messageIdString := strconv.Itoa(int(messageId))
		println("Message Id: " + messageIdString)
		switch messageId {
		case REGISTER_REQUEST_ID:
			_, err = conn.Read(buffer) // Ignoring the ID for now
			_, err = conn.Read(buffer)
			port := binary.LittleEndian.Uint32(buffer)
			portString := strconv.Itoa(int(port))
			println("Port: " + portString)
			tcpAddr, _ := net.ResolveTCPAddr("tcp", "localhost:"+portString)
			fmt.Println(conn.RemoteAddr().String())
			tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
			if err != nil {
				fmt.Println("There was a resolve error")
				return
			}
			messengers[port] = tcpConn
		}
	}
}
func getRandomNodeId(randomSource rand.Source) int {
	randGen := rand.New(randomSource)
	return randGen.Intn(9999)
}
func setupRing(messengers map[uint32]net.Conn) {
	keys := make([]uint32, 0, len(messengers))
	for k := range messengers {
		keys = append(keys, k)
		fmt.Println(strconv.Itoa(len(messengers)))
	}
	for i := 0; i < len(keys)-1; i++ {
		port, _ := strconv.Atoi(messengers[keys[i+1]].LocalAddr().String()[strings.Index(messengers[keys[i+1]].LocalAddr().String(), ":"):])
		connect := GetConnectionsDirective(uint32(port), keys[i+1])
		messengers[keys[i]].Write(GetConnectionsDirectiveBytes(connect))
	}
	port, _ := strconv.Atoi(messengers[keys[0]].LocalAddr().String()[strings.Index(messengers[keys[0]].LocalAddr().String(), ":"):])
	connect := GetConnectionsDirective(uint32(port), keys[0])
	messengers[keys[len(keys)-1]].Write(GetConnectionsDirectiveBytes(connect))
}

func commence(messengers map[uint32]net.Conn, numMessagesToSend uint32) {
	for k := range messengers {
		task := GetTaskInitiate(numMessagesToSend)
		messengers[k].Write(GetTaskInitiateBytes(task))
	}
}
