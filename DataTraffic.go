package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type DataTraffic struct {
	sourceId  uint32
	randomInt uint32
}

func GetDataTrafficBytes(traffic DataTraffic) []byte {
	//fmt.Println("Vals " + strconv.Itoa(traffic.randomInt) + " " + strconv.Itoa(traffic.sourceId))
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, DATA_TRAFFIC_ID)
	err = binary.Write(buf, binary.LittleEndian, traffic.sourceId)
	err = binary.Write(buf, binary.LittleEndian, traffic.randomInt)
	if err != nil {
		fmt.Println("Problems getting the bytes")
		return nil
	} else {
		return buf.Bytes()
	}
}
func GetDataTraffic(sourceId uint32, randomInt uint32) DataTraffic {
	return DataTraffic{sourceId: sourceId, randomInt: randomInt}
}
