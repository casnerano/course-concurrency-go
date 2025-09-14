package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/casnerano/course-concurrency-go/internal/network"
	"github.com/casnerano/course-concurrency-go/internal/network/protocol"
)

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
	fmt.Print("query > ")
	if scanner.Scan() {
		var response *protocol.Response
		response, err = client.Send(scanner.Text())
		if err != nil {
			fmt.Println("response error: ", err.Error())
			return
		}

		fmt.Printf("response: %+v\n", response)
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
