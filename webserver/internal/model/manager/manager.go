package manager

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"webserver/internal/model/httpentity"
)

const BUFF_SIZE = 1024
const DEFAULT_SERVER_PORT = uint32(8080)

type ConnectionManager interface {
	Start() error
	handleConnection(conn net.Conn)
}

type ConcurrentConnectionManger struct {
	port uint32
}

func NewConcurrentConnectionManger(port ...uint32) *ConcurrentConnectionManger {
	actualPort := DEFAULT_SERVER_PORT
	if len(port) > 0 {
		actualPort = port[0]
	}
	return &ConcurrentConnectionManger{port: actualPort}
}

func (cm *ConcurrentConnectionManger) Start() error {
	listener, err := net.Listen("tcp", "localhost:"+strconv.FormatUint(uint64(cm.port), 10))
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error in listening" + err.Error())
			continue
		}
		cm.handleConnection(conn)
	}
}

func (cm *ConcurrentConnectionManger) handleConnection(conn net.Conn) {
	defer conn.Close()

	data, err := readAll(conn)
	if err != nil {
		log.Println("Error in listening" + err.Error())
		return
	}
	log.Println("Received request\n", string(data))
	req, err := httpentity.ParseRequest(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	a := fmt.Sprintf("HTTP/1.1 200 OK\r\n\r\nRequested path:\n %s\r\n", req.Path)
	conn.Write([]byte(a))
}

func readAll(conn net.Conn) ([]byte, error) {
	data := make([]byte, 0, 4*BUFF_SIZE)
	for {
		buff := make([]byte, BUFF_SIZE)
		read, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				return data, nil
			}
			return nil, err
		}
		buff = buff[:read]
		data = append(data, buff...)
		if isReadingFinished(data) {
			break
		}
	}
	return data, nil
}

func isReadingFinished(data []byte) bool {
	len := len(data)
	if len < 4 {
		return false
	}
	return data[len-4] == '\r' && data[len-3] == '\n' && data[len-2] == '\r' && data[len-1] == '\n'
}
