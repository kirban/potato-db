package network

import (
	"bufio"
	"context"
	"github.com/kirban/potato-db/internal/network/handlers"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/kirban/potato-db/internal/config"
	"github.com/kirban/potato-db/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewTCPServer(t *testing.T) {
	tests := map[string]struct {
		logger  *zap.Logger
		config  *config.ServerConfigOptions
		handler TCPRequestHandler
		wantErr bool
	}{
		"valid config": {
			logger: createTestLogger(),
			config: &config.ServerConfigOptions{Host: "127.0.0.1", Port: 8080},
			handler: &handlers.DatabaseHandler{
				Db: createMockDatabase(),
			},
			wantErr: false,
		},
		"nil logger": {
			logger: nil,
			config: &config.ServerConfigOptions{Host: "127.0.0.1", Port: 8080},
			handler: &handlers.DatabaseHandler{
				Db: createMockDatabase(),
			},
			wantErr: true,
		},
		"nil config": {
			logger: createTestLogger(),
			config: nil,
			handler: &handlers.DatabaseHandler{
				Db: createMockDatabase(),
			},
			wantErr: true,
		},
		"nil database": {
			logger:  createTestLogger(),
			config:  &config.ServerConfigOptions{Host: "127.0.0.1", Port: 8080},
			handler: nil,
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			server, err := NewTCPServer(tt.logger, tt.config, tt.handler)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, server)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, server)
				if tt.config != nil {
					assert.Equal(t, tt.config.Host, server.host)
					assert.Equal(t, tt.config.Port, server.port)
				}
				assert.Equal(t, tt.logger, server.logger)
				assert.Equal(t, tt.handler, server.handler)
			}
		})
	}
}

func TestTCPServer_StartAndServe(t *testing.T) {
	logger := createTestLogger()
	handler := &handlers.DatabaseHandler{
		Db: createMockDatabase(),
	}
	server, err := NewTCPServer(logger, &config.ServerConfigOptions{Host: "127.0.0.1", Port: 0}, handler)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Test server startup
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.StartAndServe(ctx)
	}()

	// Wait for server to start
	time.Sleep(10 * time.Millisecond)

	// Test that server stops when context is cancelled
	select {
	case err := <-errChan:
		assert.NoError(t, err)
	case <-time.After(200 * time.Millisecond):
		t.Fatal("server did not stop within timeout")
	}
}

func TestTCPServer_Stop(t *testing.T) {
	logger := createTestLogger()
	handler := &handlers.DatabaseHandler{
		Db: createMockDatabase(),
	}
	server, err := NewTCPServer(logger, &config.ServerConfigOptions{Host: "127.0.0.1", Port: 0}, handler)
	require.NoError(t, err)

	// Test stop on nil listener
	server.Stop() // Should not panic

	// Test stop on active listener
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		server.StartAndServe(ctx)
	}()

	time.Sleep(10 * time.Millisecond)
	server.Stop()
}

func TestTCPServer_handleConnection(t *testing.T) {
	logger := createTestLogger()
	handler := &handlers.DatabaseHandler{
		Db: createMockDatabase(),
	}
	server, err := NewTCPServer(logger, &config.ServerConfigOptions{Host: "127.0.0.1", Port: 0}, handler)
	require.NoError(t, err)

	// Create a test connection pair
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	// Test normal connection handling
	go func() {
		server.handleConnection(serverConn)
	}()

	// Send test data
	testData := "test query\n"
	_, err = clientConn.Write([]byte(testData))
	require.NoError(t, err)

	// Read response
	response, err := bufio.NewReader(clientConn).ReadString('\n')
	require.NoError(t, err)
	assert.True(t, strings.HasSuffix(response, "\n"))
}

func createTestLogger() *zap.Logger {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.ErrorLevel),
		Development: false,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			MessageKey:     "msg",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := config.Build()
	return logger
}

// mockDatabase implements db.Executable for testing
type mockDatabase struct{}

func (m *mockDatabase) ExecuteQuery(q string) (string, error) {
	return "OK mock response", nil
}

func createMockDatabase() db.Executable {
	return &mockDatabase{}
}
