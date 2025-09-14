package network

import (
	"fmt"
	"net"

	"github.com/casnerano/course-concurrency-go/internal/network/protocol"
)

type ClientOptions struct {
	Address string
}

type Client struct {
	connection net.Conn
	protocol   protocol.Protocol
	options    ClientOptions
}

func NewClient(protocol protocol.Protocol, options ClientOptions) *Client {
	return &Client{
		protocol: protocol,
		options:  options,
	}
}

func (c *Client) Connect() error {
	var err error
	c.connection, err = net.Dial("tcp", c.options.Address)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}

	return nil
}

func (c *Client) Send(query string) (*protocol.Response, error) {
	request := protocol.Request{
		Payload: protocol.RequestPayload{
			RawQuery: query,
		},
	}

	if err := c.protocol.EncodeRequest(c.connection, &request); err != nil {
		return nil, fmt.Errorf("failed encode request: %w", err)
	}

	response, err := c.protocol.DecodeResponse(c.connection)
	if err != nil {
		return nil, fmt.Errorf("failed decode response: %w", err)
	}

	return response, nil
}

func (c *Client) Close() error {
	if c.connection == nil {
		return nil
	}

	return c.connection.Close()
}
