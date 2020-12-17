package tcp_socket

import (
	"fmt"
	"net"
	"testing"
)

func TestClient(t *testing.T) {

}

func Client() {

	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:9999")

	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		fmt.Println("Client connect error ! " + err.Error())
		return
	}

	defer conn.Close()

	fmt.Println(conn.LocalAddr().String() + " : Client connected!")

	onMessageReceived(conn)
}

func onMessageReceived(conn *net.TCPConn) {

	b := []byte("aaaaaaaaaaaaaaaaaaaaa\n")
	_, err := conn.Write(b)
	fmt.Println(err)
}
