package network

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/casnerano/course-concurrency-go/internal/database"
	"github.com/casnerano/course-concurrency-go/internal/logger"
	"github.com/casnerano/course-concurrency-go/internal/network/protocol"
	"github.com/casnerano/course-concurrency-go/internal/types"
)

type dbHandler interface {
	HandleQuery(ctx context.Context, req database.Query) (*types.Value, error)
}

type Server struct {
	addr      string
	listener  net.Listener
	protocol  protocol.Protocol
	dbHandler dbHandler
}

func NewServer(addr string, protocol protocol.Protocol, dbHandler dbHandler) *Server {
	return &Server{
		addr:      addr,
		protocol:  protocol,
		dbHandler: dbHandler,
	}
}

func (s *Server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	defer func() {
		err = listener.Close()
		if err != nil {
			logger.Error("failed close listener: " + err.Error())
		}
	}()

	s.listener = listener

	go func() {
		<-ctx.Done()
		if err = listener.Close(); err != nil {
			logger.Error("failed close listener: " + err.Error())
		}
	}()

	logger.Info("server started on: " + s.addr)

	for {
		conn, acceptErr := s.listener.Accept()
		if acceptErr != nil {
			if ctx.Err() != nil {
				return nil
			}

			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				continue
			}

			return fmt.Errorf("failed accept: %w", err)
		}

		go s.handleConnection(ctx, conn)
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			logger.Error("failed close connection: " + closeErr.Error())
		}
	}()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		request, decodeErr := s.protocol.DecodeRequest(reader)
		if decodeErr != nil {
			if decodeErr == io.EOF {
				return
			}

			response := protocol.Response{
				Status:       protocol.ResponseStatusCancel,
				ErrorMessage: "failed decode request: " + decodeErr.Error(),
			}

			if err := s.protocol.EncodeResponse(writer, &response); err != nil {
				logger.Error("failed encode response", "context", response)
			}

			return
		}

		dbQuery := database.Query{
			Command: request.Payload.Command,
			Key:     request.Payload.Key,
			Value:   request.Payload.Value,
		}

		value, handleQueryErr := s.dbHandler.HandleQuery(ctx, dbQuery)
		if handleQueryErr != nil {
			response := protocol.Response{
				Status:       protocol.ResponseStatusCancel,
				ErrorMessage: "failed handle query: " + handleQueryErr.Error(),
			}

			if err := s.protocol.EncodeResponse(writer, &response); err != nil {
				logger.Error("failed encode response", "context", response)

				return
			}
		}

		response := protocol.Response{
			Value:  value,
			Status: protocol.ResponseStatusOk,
		}

		if err := s.protocol.EncodeResponse(writer, &response); err != nil {
			logger.Error("failed encode response", "context", response)

			return
		}

		if err := writer.Flush(); err != nil {
			return
		}
	}
}
