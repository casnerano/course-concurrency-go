package network

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/casnerano/course-concurrency-go/internal/logger"
	"github.com/casnerano/course-concurrency-go/internal/network/protocol"
	"github.com/casnerano/course-concurrency-go/internal/types"
)

type dbHandler interface {
	HandleQuery(ctx context.Context, rawQuery string) (*types.Value, error)
}

type ServerOptions struct {
	Address        string
	MaxConnections int
	MaxMessageSize int
	IdleTimeout    time.Duration
}

type Server struct {
	listener  net.Listener
	protocol  protocol.Protocol
	dbHandler dbHandler
	options   ServerOptions

	connLimiter chan struct{}
}

func NewServer(protocol protocol.Protocol, dbHandler dbHandler, options ServerOptions) *Server {
	return &Server{
		options:   options,
		protocol:  protocol,
		dbHandler: dbHandler,

		connLimiter: make(chan struct{}, options.MaxConnections),
	}
}

func (s *Server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.options.Address)
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

	logger.Info("server started on: " + s.options.Address)

	return s.listen(ctx)
}

func (s *Server) listen(ctx context.Context) error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}

			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				continue
			}

			return fmt.Errorf("failed accept: %w", err)
		}

		s.connLimiter <- struct{}{}

		go func() {
			defer func() {
				<-s.connLimiter
			}()

			s.handleConnection(ctx, conn)
		}()
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			s.sendErrorResponse(conn, "internal error")
		}
	}()

	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			logger.Error("failed close connection: " + closeErr.Error())
		}
	}()

	for {
		readDeadline := time.Now().Add(s.options.IdleTimeout * time.Millisecond)
		if err := conn.SetReadDeadline(readDeadline); err != nil {
			s.sendErrorResponse(conn, "failed to set read deadline")
			return
		}

		request, decodeErr := s.protocol.DecodeRequest(conn)
		if decodeErr != nil {
			if decodeErr == io.EOF {
				return
			}

			s.sendErrorResponse(conn, "failed decode request")
			break
		}

		dbValue, dbHandleErr := s.dbHandler.HandleQuery(ctx, request.Payload.RawQuery)
		if dbHandleErr != nil {
			s.sendErrorResponse(conn, "failed handle query")
			continue
		}

		var payload *protocol.ResponsePayload
		if dbValue != nil {
			payload = &protocol.ResponsePayload{
				Value: dbValue,
			}
		}

		s.sendSuccessResponse(conn, payload)
	}
}

func (s *Server) sendErrorResponse(writer io.Writer, error string) {
	s.sendFailedResponse(writer, protocol.ResponseStatusError, &error)
}

func (s *Server) sendCancelResponse(writer io.Writer, error string) {
	s.sendFailedResponse(writer, protocol.ResponseStatusCancel, &error)
}

func (s *Server) sendSuccessResponse(writer io.Writer, payload *protocol.ResponsePayload) {
	s.sendResponse(writer, protocol.Response{
		Payload: payload,
		Status:  protocol.ResponseStatusOk,
	})
}

func (s *Server) sendFailedResponse(writer io.Writer, status protocol.ResponseStatus, error *string) {
	s.sendResponse(writer, protocol.Response{
		Status: status,
		Error:  error,
	})
}

func (s *Server) sendResponse(writer io.Writer, response protocol.Response) {
	if err := s.protocol.EncodeResponse(writer, &response); err != nil {
		logger.Error("failed encode response", "context", response)
	}
}
