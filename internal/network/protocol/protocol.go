package protocol

import (
	"io"

	"github.com/casnerano/course-concurrency-go/internal/types"
)

type Protocol interface {
	Send(writer io.Writer, v any) error
	Receive(reader io.Reader, v any) error
}

type Request struct {
	Payload RequestPayload `json:"payload"`
}

type RequestPayload struct {
	RawQuery string `json:"raw_query"`
}

type Response struct {
	Payload *ResponsePayload `json:"payload,omitempty"`
	Status  ResponseStatus   `json:"status"`
	Error   *string          `json:"error,omitempty"`
}

type ResponsePayload struct {
	Value *types.Value `json:"value,omitempty"`
}

type ResponseStatus string

const (
	ResponseStatusOk     ResponseStatus = "OK"
	ResponseStatusCancel ResponseStatus = "CANCEL"
	ResponseStatusError  ResponseStatus = "ERROR"
)
