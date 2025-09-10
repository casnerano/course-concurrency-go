package protocol

import (
	"encoding/json"
	"io"
)

var _ Protocol = (*JSON)(nil)

type JSON struct{}

func NewJSON() *JSON {
	return &JSON{}
}

func (j *JSON) EncodeRequest(writer io.Writer, request *Request) error {
	return json.NewEncoder(writer).Encode(request)
}

func (j *JSON) EncodeResponse(writer io.Writer, response *Response) error {
	return json.NewEncoder(writer).Encode(response)
}

func (j *JSON) DecodeRequest(reader io.Reader) (*Request, error) {
	var request Request
	if err := json.NewDecoder(reader).Decode(&request); err != nil {
		return nil, err
	}
	return &request, nil
}

func (j *JSON) DecodeResponse(reader io.Reader) (*Response, error) {
	var response Response
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}
