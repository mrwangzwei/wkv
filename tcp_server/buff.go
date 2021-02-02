package tcp_server

import (
	"bufio"
	"errors"
)

func buffReader(reader *bufio.Reader) ([]byte, error) {
	message, err := reader.ReadSlice('\n')
	if err == bufio.ErrBufferFull {
		newBuf := append([]byte{}, message...)
		for err == bufio.ErrBufferFull {
			message, err = reader.ReadBytes('\n')
			newBuf = append(newBuf, message...)
		}
		message = newBuf
	}
	if err != nil {
		return nil, err
	}
	i := len(message) - 1
	if i < 0 {
		return nil, errors.New("nil read")
	}
	return message[:i], nil
}
