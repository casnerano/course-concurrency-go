package network

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/casnerano/course-concurrency-go/internal/logger"
	"github.com/casnerano/course-concurrency-go/internal/network/protocol"
)

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

		request, decodeErr := s.protocol.DecodeRequest(conn, s.options.MaxMessageSize)
		if decodeErr != nil {
			if decodeErr == io.EOF {
				return
			}

			s.sendErrorResponse(conn, "failed decode request")
			continue
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
