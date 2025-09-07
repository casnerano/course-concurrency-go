package protocol

import (
	"encoding/json"
	"fmt"
	"io"
)

func RequestDecode(reader io.Reader) (*Request, error) {
	request := Request{}
	if err := json.NewDecoder(reader).Decode(&request); err != nil {
		return nil, fmt.Errorf("failed request decode: %w", err)
	}

	return &request, nil
}

func ResponseDecode(reader io.Reader) (*Response, error) {
	response := Response{}
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed response decode: %w", err)
	}

	return &response, nil
}
