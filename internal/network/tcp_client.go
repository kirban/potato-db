package network

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var defaultBufferSize = 4 << 10

type TCPClient struct {
	conn        net.Conn
	idleTimeout time.Duration
	bufferSize  int
}

func NewTCPClient(address string, idleTimeout time.Duration, bufferSize int) (*TCPClient, error) {
	conn, err := net.Dial("tcp", address)

	if err != nil {
		return nil, err
	}

	if bufferSize == 0 {
		bufferSize = defaultBufferSize
	}

	client := &TCPClient{
		conn:        conn,
		idleTimeout: idleTimeout,
		bufferSize:  bufferSize,
	}

	if client.idleTimeout != 0 {
		if err := conn.SetDeadline(time.Now().Add(client.idleTimeout)); err != nil {
			return nil, fmt.Errorf("failed to set deadline for connection: %w", err)
		}
	}

	return client, nil
}

func (c *TCPClient) Send(request []byte) ([]byte, error) {
	if _, err := c.conn.Write(request); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	response := make([]byte, c.bufferSize)
	count, err := c.conn.Read(response)

	if err != nil && err != io.EOF {
		return nil, err
	} else if count == c.bufferSize {
		return nil, errors.New("small buffer size")
	}

	return response[:count], nil
}

func (c *TCPClient) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
