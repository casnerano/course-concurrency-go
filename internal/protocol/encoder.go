package protocol

import (
	"encoding/json"
	"fmt"
)

func RequestEncode(request Request) ([]byte, error) {
	bRequest, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed request encode")
	}

	return bRequest, nil
}

func ResponseEncode(response Response) ([]byte, error) {
	bResponse, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed response encode")
	}

	return bResponse, nil
}
