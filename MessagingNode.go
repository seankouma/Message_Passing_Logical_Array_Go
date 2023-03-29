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
	"sync/atomic"
)

func runMessageNode(host string, port string, destination_host string, destination_port string, messageCount string) {
	id := getRandomId()
	var registrySocket net.Conn
	go registerNode(id, port, destination_port, registrySocket)
	numData := new(uint32)
	*numData = 0
	peerSocket := new(net.Conn)
	go handleRequest(host, port, numData, peerSocket)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ready for input")
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading")
			return
		}
		if strings.Contains(input, "status") {
			fmt.Println("Num Messages " + strconv.Itoa(int(*numData)))
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

func handleRequest(host string, port string, numData *uint32, peerSocket *net.Conn) {
	fmt.Println("Handle Request")
	listen := getListener(host, port)
	for {
		conn := getConnection(listen)
		go listenForMessages(conn, port, numData, peerSocket)
	}
}
func listenForMessages(conn net.Conn, port string, numData *uint32, peerSocket *net.Conn) {
	for {
		messageId := getMessageId(conn)
		//messageIdString := strconv.Itoa(int(messageId))
		//println("Message Id: " + messageIdString)
		switch messageId {
		case CONNECTIONS_DIRECTIVE_ID:
			*peerSocket = handleConnectionsDirective(conn)
			if peerSocket != nil {
				fmt.Println("PeerSocket is not null")
			}
		case DATA_TRAFFIC_ID:
			handleDataTraffic(conn, peerSocket, port, numData)
			atomic.AddUint32(numData, 1)
		case TASK_INITIATE_ID:
			fmt.Println("Task Initiate")
			handleTaskInitiate(conn, peerSocket, port)
		default:
			fmt.Println("Something went wrong: " + strconv.Itoa(int(messageId)))
		}
	}
}

func getMessageId(conn net.Conn) uint32 {
	buffer := make([]byte, 4)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	return binary.LittleEndian.Uint32(buffer)
}
func getConnection(listen net.Listener) net.Conn {
	conn, err := listen.Accept()
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
func getListener(host string, port string) net.Listener {
	listen, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	return listen
}
func handleConnectionsDirective(conn net.Conn) net.Conn {
	buffer := make([]byte, 4)
	conn.Read(buffer)
	nodeID := binary.LittleEndian.Uint32(buffer)
	nodeIDStr := strconv.Itoa(int(nodeID))
	println("Node ID: " + nodeIDStr)
	conn.Read(buffer)
	port := binary.LittleEndian.Uint32(buffer)
	portStr := strconv.Itoa(int(port))
	println("Port String: " + portStr)
	tcpServer, err := net.ResolveTCPAddr("tcp", "localhost:"+portStr)
	if err != nil {
		log.Fatal(err)
	}
	peerSocket, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		log.Fatal(err)
	}
	return peerSocket
}
func handleDataTraffic(conn net.Conn, peerSocket *net.Conn, port string, numData *uint32) {
	buffer := make([]byte, 4)
	conn.Read(buffer)
	nodeId := binary.LittleEndian.Uint32(buffer)
	nodeIdNum := strconv.Itoa(int(nodeId))
	conn.Read(buffer)
	randomVal := binary.LittleEndian.Uint32(buffer)
	if atomic.LoadUint32(numData)%10000 == 0 {
		println("Messages received: ", atomic.LoadUint32(numData))
	}
	if nodeIdNum != port && peerSocket != nil {
		dataTraffic := GetDataTraffic(nodeId, randomVal)
		if peerSocket == nil {
			fmt.Println("PeerSocket is Null")
		}
		forwardMessage(peerSocket, dataTraffic)
	}
}

func forwardMessage(peerSocket *net.Conn, dataTraffic DataTraffic) {
	_, err := (*peerSocket).Write(GetDataTrafficBytes(dataTraffic))
	if err != nil {
		log.Fatal(err)
	}
}
func handleTaskInitiate(conn net.Conn, peerSocket *net.Conn, port string) {
	buffer := make([]byte, 4)
	conn.Read(buffer)
	numMessages := binary.LittleEndian.Uint32(buffer)
	fmt.Println("Num of messages: " + strconv.Itoa(int(numMessages)))
	go sendDataTraffic(peerSocket, numMessages, port)
}
func sendDataTraffic(peerSocket *net.Conn, numMessages uint32, port string) {
	portInt, _ := strconv.Atoi(port)
	randomSource := rand.NewSource(35)
	randGen := rand.New(randomSource)
	for i := 0; i < int(numMessages); i++ {
		dataTraffic := GetDataTraffic(uint32(portInt), randGen.Uint32())
		_, err := (*peerSocket).Write(GetDataTrafficBytes(dataTraffic))
		if err != nil {
			log.Fatal(err)
		}
	}
}
