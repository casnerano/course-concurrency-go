package protocol

import (
	"io"

	"github.com/casnerano/course-concurrency-go/internal/types"
)

type Protocol interface {
	EncodeRequest(writer io.Writer, request *Request) error
	EncodeResponse(writer io.Writer, response *Response) error
	DecodeRequest(reader io.Reader, maxSize int) (*Request, error)
	DecodeResponse(reader io.Reader) (*Response, error)
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
