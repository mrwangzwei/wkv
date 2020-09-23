package dns

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

func simpleSend(domain string) ([]string, error) {
	ns, err := net.LookupHost(domain)
	if err != nil {
		return nil, err
	}
	ipArr := make([]string, len(ns))
	for index, n := range ns {
		ipArr[index] = n
	}
	return ipArr, nil
}

func Send(addr, domain string) {
	conn, err := net.Dial("udp", addr)
	defer conn.Close()

	requestHeader := dnsHeader{
		Id:      0x0010,
		Qdcount: 1,
		Ancount: 0,
		Nscount: 0,
		Arcount: 0,
	}
	requestHeader.SetFlag(0, 0, 0, 0, 1, 0, 0)

	requestQuery := dnsQuery{
		QuestionType:  1,
		QuestionClass: 1,
	}

	var buffer bytes.Buffer
	err = binary.Write(&buffer, binary.BigEndian, requestHeader)
	err = binary.Write(&buffer, binary.BigEndian, ParseDomainName(domain))
	err = binary.Write(&buffer, binary.BigEndian, requestQuery)

	buf := make([]byte, dnsPacketLen)
	t1 := time.Now()
	_, err = conn.Write(buffer.Bytes())
	fmt.Println("conn.Write", err)

	length, err := conn.Read(buf)
	t := time.Now().Sub(t1)
	fmt.Println("conn.Read", string(buf))
	fmt.Println("conn.Read", length, t)

}

//read 每个包的大小固定是512，其中第一字节做标志位
//bagMap 作为收包存储器，按包的顺序存储到对应的索引下
func SendJson(addr, domain string) {
	mmp := make(map[string]string)
	mmp["domain"] = domain
	val, _ := json.Marshal(mmp)
	conn, err := net.Dial("udp", addr)
	defer conn.Close()

	t1 := time.Now()
	_, err = conn.Write(val)
	fmt.Println("conn.Write", err)

	bagMap := make([][]byte, MaxBufCount)
	for {
		buf := make([]byte, 512)
		_, err = conn.Read(buf)

		if err != nil {
			fmt.Println(err)
			break
		}
		if buf[0] < BufNextExist {
			bagMap[buf[0]] = buf
			break
		} else {
			bagMap[buf[0]-BufNextExist] = buf
		}
	}

	var result []byte
	lostCount := 0
	for _, item := range bagMap {
		if len(item) > 1 {
			result = append(result, item[1:]...)
		} else {
			lostCount++
			if lostCount >= 3 {
				log.Printf("nil bag over %d times", lostCount)
				break
			}
		}
	}

	t := time.Now().Sub(t1)
	fmt.Println(string(result), t)
}
