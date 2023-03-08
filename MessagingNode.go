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

func runMessageNode(host string, port string, destination_host string, destination_port string, messageCount string) {
	id := getRandomId()
	var peerSocket net.Conn
	var registrySocket net.Conn
	go registerNode(id, port, destination_port, registrySocket)
	go handleRequest(host, port, peerSocket)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ready for input")
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading")
			return
		}
		if strings.Contains(input, "connect") {
			fmt.Println("Read input")
			go sendMessages(destination_host, destination_port, messageCount, id)
		}
	}
}

func getRandomId() int {
	randomSource := rand.NewSource(42)
	randGen := rand.New(randomSource)
	return randGen.Intn(9999)
}

func registerNode(id int, port string, destination_port string, registrySocket net.Conn) {
	tcpServer, err := net.ResolveTCPAddr("tcp", "localhost:"+destination_port)
	if err != nil {
		fmt.Println("There was a resolve error")
		return
	}

	registrySocket, err = net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		fmt.Println("There was a dial error")
		return
	}
	portInt, err := strconv.Atoi(port)
	register := GetRegisterRequest(uint32(id), uint32(portInt))
	_, err = registrySocket.Write(GetRegisterRequestBytes(register))
	if err != nil {
		fmt.Println("There was a writing error")
		return
	} else {
		fmt.Println("Wrote registration request")
	}
}

func handleRequest(host string, port string, peerSocket net.Conn) {
	fmt.Println("Handle Request")
	listen, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// close listener
	defer listen.Close()
	conn, err := listen.Accept()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Println("Made a connection")
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
		case DATA_TRAFFIC_ID:
			_, err = conn.Read(buffer)
			randomVal := binary.LittleEndian.Uint32(buffer)
			randomNum := strconv.Itoa(int(randomVal))
			println("Random num: " + randomNum)
		case CONNECTIONS_DIRECTIVE_ID:
			_, err = conn.Read(buffer)
			nodeID := binary.LittleEndian.Uint32(buffer)
			nodeIDStr := strconv.Itoa(int(nodeID))
			println("Node ID: " + nodeIDStr)
			_, err = conn.Read(buffer)
			port := binary.LittleEndian.Uint32(buffer)
			portStr := strconv.Itoa(int(port))
			println("Port String: " + portStr)
			tcpServer, err := net.ResolveTCPAddr("tcp", "localhost:"+portStr)
			if err != nil {
				fmt.Println("There was a resolve error")
				return
			}

			peerSocket, err = net.DialTCP("tcp", nil, tcpServer)
			if err != nil {
				fmt.Println("There was a dial error")
				return
			}
		}

	}
	// close conn
	conn.Close()
}

func sendMessages(destinationHost string, destinationPort string, messageCount string, id int) {
	fmt.Println("Send Messages")
	tcpServer, err := net.ResolveTCPAddr("tcp", destinationHost+":"+destinationPort)
	if err != nil {
		fmt.Println("There was a resolve error")
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		fmt.Println("There was a dial error")
		return
	}

	messageCountInt, _ := strconv.Atoi(messageCount)
	for i := 0; i < messageCountInt; i++ {
		if err != nil {
			log.Fatal(err)
		}
		randomSource := rand.NewSource(42)
		randGen := rand.New(randomSource)
		dataTraffic := GetDataTraffic(uint32(id), randGen.Uint32())
		_, err = conn.Write(GetDataTrafficBytes(dataTraffic))
		if err != nil {
			fmt.Println("There was a writing error")
			return
		} else {
			fmt.Println("Wrote the bytes")
		}
	}
}
