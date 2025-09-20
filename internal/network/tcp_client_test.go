package network

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"syscall"
	"testing"
	"time"
)

func TestTCPClient(t *testing.T) {
	t.Parallel()

	const serverResponse = "hello from server"
	const serverAddress = "localhost:11111"
	listener, err := net.Listen("tcp", serverAddress)
	require.NoError(t, err)

	go func() {
		for {
			conn, err := listener.Accept()

			if err != nil {
				return
			}

			_, err = conn.Read(make([]byte, 2048))
			require.NoError(t, err)

			_, err = conn.Write([]byte(serverResponse))
			require.NoError(t, err)
		}
	}()

	tests := map[string]struct {
		request string
		client  func() *TCPClient

		expectedResponse string
		expectedErr      error
	}{
		"client with incorrect server address": {
			request: "hello server",
			client: func() *TCPClient {
				client, err := NewTCPClient("localhost:1010", 0, 0)
				require.ErrorIs(t, err, syscall.ECONNREFUSED)
				return client
			},
			expectedResponse: serverResponse,
		},
		"client with small max message size": {
			request: "hi there",
			client: func() *TCPClient {
				client, err := NewTCPClient("localhost:11111", 0, 5)
				require.NoError(t, err)
				return client
			},
			expectedErr: errors.New("small buffer size"),
		},
		"client with idle timeout": {
			request: "hello from client",
			client: func() *TCPClient {
				client, err := NewTCPClient("localhost:11111", 2*time.Second, 0)
				require.NoError(t, err)
				return client
			},
			expectedResponse: serverResponse,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := test.client()
			if client == nil {
				return
			}

			response, err := client.Send([]byte(test.request))
			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedResponse, string(response))
			client.Close()
		})
	}
}
