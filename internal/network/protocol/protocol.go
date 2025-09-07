package protocol

import (
	"io"

	"github.com/casnerano/course-concurrency-go/internal/types"
)

type Protocol interface {
	EncodeRequest(writer io.Writer, request *Request) error
	EncodeResponse(writer io.Writer, response *Response) error
	DecodeRequest(reader io.Reader) (*Request, error)
	DecodeResponse(reader io.Reader) (*Response, error)
}

type Request struct {
	Payload RequestPayload
}

type RequestPayload struct {
	Command types.Command
	Key     *types.Key
	Value   *types.Value
}

type Response struct {
	Value *types.Value

	Status       ResponseStatus
	ErrorMessage string
}

type ResponseStatus string

const (
	ResponseStatusOk     ResponseStatus = "OK"
	ResponseStatusCancel ResponseStatus = "CANCEL"
)
