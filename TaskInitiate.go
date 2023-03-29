package main

import (
    "bytes"
    "encoding/binary"
    "fmt"
)

type TaskInitiate struct {
    quantity uint32
}

func GetTaskInitiateBytes(task TaskInitiate) []byte {
    buf := new(bytes.Buffer)
    err := binary.Write(buf, binary.LittleEndian, TASK_INITIATE_ID)
    err = binary.Write(buf, binary.LittleEndian, task.quantity)
    if err != nil {
        fmt.Println("Problems getting the bytes")
        return nil
    } else {
        return buf.Bytes()
    }
}
func GetTaskInitiate(quantity uint32) TaskInitiate {
    return TaskInitiate{quantity: quantity}
}
