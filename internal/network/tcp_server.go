package network

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/kirban/potato-db/internal/config"
	"github.com/kirban/potato-db/internal/db"
	"go.uber.org/zap"
	"net"
	"sync"
	"time"
)

type TCPRouteHandler func(string) string

type TCPServer struct {
	host           string
	port           int
	logger         *zap.Logger
	listener       net.Listener
	database       db.Executable
	bufferSize     int
	idleTimeout    time.Duration
	maxConnections int
	semaphore      chan struct{}
}

func NewTCPServer(logger *zap.Logger, config *config.ServerConfigOptions, database db.Executable) (*TCPServer, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	if config == nil {
		return nil, errors.New("config is invalid")
	}

	if database == nil {
		return nil, errors.New("database is invalid")
	}

	return &TCPServer{
		host:           config.Host,
		port:           config.Port,
		logger:         logger,
		database:       database,
		bufferSize:     config.BufferSize,
		maxConnections: config.MaxConnections,
		semaphore:      make(chan struct{}, config.MaxConnections),
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

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				// net.ErrClosed when listener closed on shutdown
				if errors.Is(err, net.ErrClosed) {
					s.logger.Info("listener closed")
					return
				}
				s.logger.Error("accept error", zap.Error(err))
				continue
			}

			select {
			case s.semaphore <- struct{}{}:
				go s.handleConnection(conn)
			default:
				s.logger.Warn("too many connections", zap.String("remote", conn.RemoteAddr().String()))
				_ = conn.Close()
			}

		}
	}()

	<-ctx.Done()
	s.Stop()
	wg.Wait()

	return nil
}

func (s *TCPServer) Stop() {
	if s.listener != nil {
		_ = s.listener.Close()
	}
	s.logger.Info("server stopped")
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	s.logger.Info("client connected", zap.String("remote", conn.RemoteAddr().String()))

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
		<-s.semaphore
	}(conn)

	request := make([]byte, 0, s.bufferSize)
	scanner := bufio.NewScanner(conn)
	scanner.Buffer(request, s.bufferSize)

	for scanner.Scan() {
		line := scanner.Text()
		s.logger.Info("received", zap.String("msg", line), zap.String("remote", conn.RemoteAddr().String()))

		response, err := s.database.ExecuteQuery(line)
		if err != nil {
			s.logger.Error("database query failed", zap.Error(err))
			response = "ERROR database query failed"
		}

		responseStr := fmt.Sprintf("%v", response)
		if _, err := conn.Write([]byte(responseStr + "\n")); err != nil {
			s.logger.Error("failed to write response", zap.Error(err))
			return
		}
	}
}
