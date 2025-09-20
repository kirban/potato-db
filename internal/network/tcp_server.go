package network

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/kirban/potato-db/internal/config"
	"github.com/kirban/potato-db/internal/network/handlers"
	"go.uber.org/zap"
)

type TCPServer struct {
	host     string
	port     int
	logger   *zap.Logger
	listener net.Listener
	handler  func(string) string
}

func NewTCPServer(logger *zap.Logger, config *config.ServerConfigOptions) (*TCPServer, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	if config == nil {
		return nil, errors.New("config is invalid")
	}

	return &TCPServer{
		host:     config.Host,
		port:     config.Port,
		logger:   logger,
		listener: nil,
		handler:  handlers.HandleQuery,
	}, nil
}

func (s *TCPServer) StartAndServe(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("accept loop panic", zap.Any("recovery", r), zap.Stack("stack"))
			s.Stop()
		}
	}()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))

	if err != nil {
		return fmt.Errorf("failed to start tcp server %v", err)
	}

	s.listener = listener

	s.logger.Info(fmt.Sprintf("TCP-server started at %s:%d", s.host, s.port))

	go func() {
		<-ctx.Done()
		s.logger.Info("stopping tcp server (ctx done)")
		_ = s.listener.Close()
	}()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			// net.ErrClosed when listener closed on shutdown
			if errors.Is(err, net.ErrClosed) {
				s.logger.Info("listener closed")
				return nil
			}
			s.logger.Error("accept error", zap.Error(err))
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *TCPServer) Stop() {
	if s.listener != nil {
		_ = s.listener.Close()
	}
	s.logger.Info("server stopped")
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	s.logger.Info("client connected")

	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("connection goroutine panic", zap.Any("recovery", r), zap.Stack("stack"))
		}
	}()

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
