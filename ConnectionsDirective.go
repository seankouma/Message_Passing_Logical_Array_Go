package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type ConnectionsDirective struct {
	id   uint32
	port uint32
}

func GetConnectionsDirectiveBytes(connect ConnectionsDirective) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, CONNECTIONS_DIRECTIVE_ID)
	err = binary.Write(buf, binary.LittleEndian, connect.id)
	err = binary.Write(buf, binary.LittleEndian, connect.port)
	if err != nil {
		fmt.Println("Problems getting the bytes")
		return nil
	} else {
		return buf.Bytes()
	}
}
func GetConnectionsDirective(id uint32, port uint32) ConnectionsDirective {
	return ConnectionsDirective{id: id, port: port}
}
