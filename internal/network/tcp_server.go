package network

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/kirban/potato-db/internal/config"
	"github.com/kirban/potato-db/internal/network/handlers"
	"go.uber.org/zap"
	"net"
)

type TCPServer struct {
	host         string
	port         int
	logger       *zap.Logger
	listener     net.Listener
	handler      func(string) string
	isRunning    bool
	shutdownChan chan struct{}
}

func NewTCPServer(logger *zap.Logger, config *config.ServerConfigOptions) (*TCPServer, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &TCPServer{
		host:         config.Host,
		port:         config.Port,
		logger:       logger,
		listener:     nil,
		handler:      handlers.HandleQuery,
		isRunning:    false,
		shutdownChan: make(chan struct{}),
	}, nil
}

func (s *TCPServer) StartAndServe() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))

	if err != nil {
		return fmt.Errorf("failed to start tcp server %v", err)
	}

	s.isRunning = true
	s.listener = listener

	fmt.Printf("TCP-server started at %s:%d \n", s.host, s.port)

	for s.isRunning {
		select {
		case <-s.shutdownChan:
			s.isRunning = false
			s.logger.Info("shutting down tcp server")
			return nil
		default:
			conn, err := s.listener.Accept()

			if err != nil {
				if !s.isRunning {
					s.logger.Info("server stopped", zap.Error(err))
					return nil
				}
				s.logger.Error("accept error", zap.Error(err))
				continue
			}

			go s.handleConnection(conn)
		}
	}

	return nil
}

func (s *TCPServer) Stop() {
	s.isRunning = false
	close(s.shutdownChan)
	s.listener.Close()

	s.logger.Info("server stopped")
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	s.logger.Info("client connected")

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			s.logger.Error("failed to close connection", zap.Error(err))
		}
	}(conn)

	data, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		s.logger.Error("failed to read request data", zap.Error(err))
		return
	}

	request := data[:len(data)-1]
	response := s.handler(request)
	if _, err := conn.Write([]byte(response + "\n")); err != nil {
		s.logger.Error("failed to write response", zap.Error(err))
		return
	}
}
