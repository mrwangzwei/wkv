package dns

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
)

//dns数据包长度
const dnsPacketLen = 512

type dnsHeader struct {
	Id, Bits, Qdcount, Ancount, Nscount, Arcount uint16
}

func (header *dnsHeader) SetFlag(QR uint16, OperationCode uint16, AuthoritativeAnswer uint16, Truncation uint16, RecursionDesired uint16, RecursionAvailable uint16, ResponseCode uint16) {
	header.Bits = QR<<15 + OperationCode<<11 + AuthoritativeAnswer<<10 + Truncation<<9 + RecursionDesired<<8 + RecursionAvailable<<7 + ResponseCode
}

type dnsQuery struct {
	QuestionType  uint16
	QuestionClass uint16
}

func ParseDomainName(domain string) []byte {
	var (
		buffer   bytes.Buffer
		segments []string = strings.Split(domain, ".")
	)
	for _, seg := range segments {
		binary.Write(&buffer, binary.BigEndian, byte(len(seg)))
		binary.Write(&buffer, binary.BigEndian, []byte(seg))
	}
	binary.Write(&buffer, binary.BigEndian, byte(0x00))

	return buffer.Bytes()
}

func unpack(rdio.Reader) ([]*Answer, error) {
	var (
		reader = bufio.NewReader(rd)
		data   []byte // 应答数据包缓存
		buf    []byte // 临时缓存
		err    error
		n      int
	)
	// 拆Header
	// ...
	// 拆Question
	question := new(Question)
	if buf, err = reader.ReadBytes(eof); err != nil { // 域名以0x00结尾
		return nil, err
	}
	data = append(data, buf...)
	question.QName = buf
	buf = make([]byte, 4)
	if n, err = reader.Read(buf); err != nil || n < 4 {
		return nil, err
	}
	data = append(data, buf...)

	binary.Read(bytes.NewBuffer(buf[0:2]), binary.BigEndian, &question.QType)
	binary.Read(bytes.NewBuffer(buf[2:]), binary.BigEndian, &question.QClass)
	// 拆Answer(s)
	answers := make([]*Answer, header.ANCount)
	buf, _ = reader.Peek(59)
	for i := 0; i < int(header.ANCount); i++ { // 根据Header中的ANCOUNT标识判断有几个Answer
		answer := new(Answer)
		// NAME
		var b byte
		var p uint16
		for {
			if b, err = reader.ReadByte(); err != nil {
				return nil, err
			}
			data = append(data, b)
			if b&pointer == pointer { // pointer是一个值为0xC0的byte类型常量
				buf = []byte{b ^ pointer, 0}
				if b, err = reader.ReadByte(); err != nil {
					return nil, err
				}
				data = append(data, b)
				buf[1] = b
				binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &p)
				if buf = getRefData(data, p); len(buf) == 0 {
					return nil, errors.New("invalid answer packet")
				}
				answer.Name = append(answer.Name, buf...)
				break
			} else {
				answer.Name = append(answer.Name, b)
				if b == eof {
					break
				}
			}
		}
		// TYPE、CLASS、TLL、RDLENGTH等其他数据
		// ...
		// RDATA
		buf = make([]byte, int(answer.RDLength))
		if n, err = reader.Read(buf); err != nil || n < int(answer.RDLength) {
			return nil, err
		}
		data = append(data, buf...)
		// 调用之前定义的SetRData()函数处理不同类型的RDATA
		if err = answer.SetRData(buf, data); err != nil {
			return nil, err
		}
		answers[i] = answer
	}
	// 拆Authority和Additional，如果有的话
	return answers, nil
}
