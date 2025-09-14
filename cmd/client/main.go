package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/casnerano/course-concurrency-go/internal/network"
	"github.com/casnerano/course-concurrency-go/internal/network/protocol"
)

const promptPrefix = "query > "

var address = flag.String("address", "", "server address")

func main() {
	if address == nil || *address == "" {
		log.Fatal("empty server address")
	}

	client, err := buildClient(*address)
	if err != nil {
		log.Fatal("failed build client: ", err.Error())
	}

	if err = client.Connect(); err != nil {
		log.Fatal("failed connect to server")
	}
	defer func() {
		if err = client.Close(); err != nil {
			log.Println("failed close connection: ", err.Error())
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(promptPrefix)

		if !scanner.Scan() {
			break
		}

		var response *protocol.Response
		response, err = client.Send(scanner.Text())
		if err != nil {
			fmt.Println("response error: ", err.Error())
			continue
		}

		switch response.Status {
		case protocol.ResponseStatusOk:
			if response.Payload != nil {
				fmt.Printf("value: %v\n", response.Payload.Value)
			} else {
				fmt.Println("Accepted")
			}
		case protocol.ResponseStatusError, protocol.ResponseStatusCancel:
			var errText any
			if response.Error != nil {
				errText = *response.Error
			}
			fmt.Printf("%s: %v\n", strings.ToLower(string(response.Status)), errText)
		default:
			fmt.Println("unknown response status:", response.Status)
		}
	}
}

func init() {
	flag.Parse()
}

func buildClient(address string) (*network.Client, error) {
	networkOptions := network.ClientOptions{
		Address: address,
	}

	client := network.NewClient(
		protocol.NewJSON(),
		networkOptions,
	)

	return client, nil
}
