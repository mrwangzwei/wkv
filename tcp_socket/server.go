package tcp_socket

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"testing"
)

func TestTcpServer(t *testing.T) {
	TcpServer()
}

func TcpServer() {
	var tcpAddr *net.TCPAddr
	//通过ResolveTCPAddr实例一个具体的tcp断点
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:9999")
	//打开一个tcp断点监听
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	defer tcpListener.Close()
	fmt.Println("Server ready to read ...")
	//循环接收客户端的连接，创建一个协程具体去处理连接
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("A client connected :" + tcpConn.RemoteAddr().String())
		go handleConn(tcpConn)
	}
}

func handleConn(conn *net.TCPConn) {
	//tcp连接的地址
	ipStr := conn.RemoteAddr().String()

	defer func() {
		fmt.Println(" Disconnected : " + ipStr)
		conn.Close()
	}()

	//获取一个连接的reader读取流
	reader := bufio.NewReaderSize(conn, 20)
	i := 0
	//接收并返回消息
	for {
		i++
		message, err := reader.ReadSlice('\n')
		if err != nil || err == io.EOF {
			fmt.Println("ReadSlice err", err)
			break
		}
		fmt.Println("ReadSlice", string(message))
		if i > 10 {
			break
		}
	}
}
