package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func Send(addr, domain string) {
	var err error
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

	buf := make([]byte, 512)
	t1 := time.Now()
	_, err = conn.Write(buffer.Bytes())
	fmt.Println("conn.Write", err)

	length, err := conn.Read(buf)
	t := time.Now().Sub(t1)
	fmt.Println("conn.Read", buf, length, t)

}
