package protocol

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

const defaultMaxBufferSize = 4096

type options struct {
	maxBufferSize int
}

type Option func(option *options)

var _ Protocol = (*JSON)(nil)

type JSON struct {
	maxBufferSize int
}

func WithMaxBufferSize(bufferSize int) Option {
	return func(option *options) {
		option.maxBufferSize = bufferSize
	}
}

func NewJSON(opts ...Option) *JSON {
	localOptions := options{
		maxBufferSize: defaultMaxBufferSize,
	}

	for _, opt := range opts {
		opt(&localOptions)
	}

	return &JSON{
		maxBufferSize: localOptions.maxBufferSize,
	}
}

func (j *JSON) Send(writer io.Writer, v any) error {
	return j.encode(writer, v)
}

func (j *JSON) Receive(reader io.Reader, v any) error {
	bufReader := bufio.NewReader(reader)

	bMessage, err := bufReader.ReadBytes('\n')
	if err != nil {
		return fmt.Errorf("read error: %w", err)
	}

	if len(bMessage) > 0 && bMessage[len(bMessage)-1] == '\n' {
		bMessage = bMessage[:len(bMessage)-1]
	}

	if len(bMessage) > j.maxBufferSize {
		return fmt.Errorf("message too large: %d bytes", len(bMessage))
	}

	return j.decode(bytes.NewReader(bMessage), v)
}

func (j *JSON) encode(writer io.Writer, v any) error {
	return json.NewEncoder(writer).Encode(v)
}

func (j *JSON) decode(reader io.Reader, v any) error {
	return json.NewDecoder(reader).Decode(v)
}
