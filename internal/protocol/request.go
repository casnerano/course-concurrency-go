package protocol

import "github.com/casnerano/course-concurrency-go/internal/types"

type RequestPayload struct {
	Command types.Command
	Key     *types.Key
	Value   *types.Value
}

type Request struct {
	Payload RequestPayload
}
