package protocol

import "github.com/casnerano/course-concurrency-go/internal/types"

type ResponseStatus string

const (
	ResponseStatusOk     ResponseStatus = "OK"
	ResponseStatusCancel ResponseStatus = "CANCEL"
)

type Response struct {
	Value *types.Value

	Status ResponseStatus
	Error  error
}
