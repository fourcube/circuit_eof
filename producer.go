package main

import (
	"log"
	"os"

	circuit "github.com/gocircuit/circuit/client"
)

/*
Producer ...
*/
type Producer struct {
	client *circuit.Client
}

/*
NewProducer ...
*/
func NewProducer() *Producer {
	circuitAddress := os.Getenv("CIRCUIT_ADDRESS")
	if !isSet(circuitAddress) {
		circuitAddress = "228.8.8.8:8822"
	}

	client := circuit.DialDiscover(circuitAddress, nil)

	return &Producer{client}
}

/*
Produce "data" on anchor "channel"
*/
func (p *Producer) Produce(data []byte) {
	tokenAnchor := getAnchor(p.client)

	var err error
	var channel circuit.Chan

	// Check if there already is a channle
	existingChannel := tokenAnchor.Get()
	if existingChannel != nil {
		channel = existingChannel.(circuit.Chan)
	}

	if channel == nil {
		channel, err = tokenAnchor.MakeChan(1)

		if err != nil {
			log.Printf("Error while creating token anchor %v", err)
			return
		}
	}

	writer, err := channel.Send()

	if err != nil {
		log.Printf("Error while opening channel for writing %v", err)
		return
	}

	writer.Write(data)
	writer.Close()
}

func getAnchor(client *circuit.Client) (anchor circuit.Anchor) {
	// Get the first available server
	for _, v := range client.View() {
		anchor = v
		break
	}

	anchor = anchor.Walk([]string{"channel"})
	return
}

func isSet(v string) bool {
	return len(v) != 0
}
