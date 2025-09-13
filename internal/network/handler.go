package network

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/casnerano/course-concurrency-go/internal/logger"
	"github.com/casnerano/course-concurrency-go/internal/network/protocol"
)

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			s.sendErrorResponse(conn, fmt.Sprintf("internal error: %s", err))
		}
	}()

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
