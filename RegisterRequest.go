package main

import (
    "bytes"
    "encoding/binary"
    "fmt"
)

type RegisterRequest struct {
    id   uint32
    port uint32
}

func GetRegisterRequestBytes(register RegisterRequest) []byte {
    buf := new(bytes.Buffer)
    err := binary.Write(buf, binary.LittleEndian, REGISTER_REQUEST_ID)
    err = binary.Write(buf, binary.LittleEndian, register.id)
    err = binary.Write(buf, binary.LittleEndian, register.port)
    if err != nil {
        fmt.Println("Problems getting the bytes")
        return nil
    } else {
        return buf.Bytes()
    }
}
func GetRegisterRequest(id uint32, port uint32) RegisterRequest {
    return RegisterRequest{id: id, port: port}
}
