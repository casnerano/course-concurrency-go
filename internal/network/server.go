package network

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/casnerano/course-concurrency-go/internal/logger"
	"github.com/casnerano/course-concurrency-go/internal/network/protocol"
	"github.com/casnerano/course-concurrency-go/internal/types"
)

type dbHandler interface {
	HandleQuery(ctx context.Context, raqQuery string) (*types.Value, error)
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

	for {
		request, decodeErr := s.protocol.DecodeRequest(conn)
		if decodeErr != nil {
			if decodeErr == io.EOF {
				return
			}

			s.sendErrorResponse(conn, "failed decode request: "+decodeErr.Error())
			continue
		}

		dbValue, dbHandleErr := s.dbHandler.HandleQuery(ctx, request.Payload.RawQuery)
		if dbHandleErr != nil {
			s.sendErrorResponse(conn, "failed handle query: "+dbHandleErr.Error())
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
