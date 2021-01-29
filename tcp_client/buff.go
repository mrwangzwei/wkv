package tcp_client

import (
	"bufio"
	"errors"
)

func buffReader(reader *bufio.Reader) ([]byte, error) {
	message, err := reader.ReadSlice('\n')
	if err == bufio.ErrBufferFull {
		newBuf := append([]byte{}, message...)
		for err == bufio.ErrBufferFull {
			message, err = reader.ReadSlice('\n')
			newBuf = append(newBuf, message...)
		}
		message = newBuf
	}
	if err != nil {
		return nil, err
	}
	i := len(message) - 2
	if i < 0 {
		return nil, errors.New("nil read")
	}
	return message[:i], nil
}
