package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"

	circuit "github.com/gocircuit/circuit/client"
)

func TestProducesData(t *testing.T) {
	stop := make(chan struct{})
	go runServer("228.8.8.8:8822", stop)
	defer stopServer(stop)

	p := NewProducer()

	// This should be sent across the channel
	want := []byte("foobar")
	p.Produce([]byte("foobar"))

	// Read the data from the channel
	// err will be EOF
	data, err := readData(p.client)

	// err == EOF is false?!
	if err != nil && err != io.EOF {
		t.Errorf("Read data failed: %v, data: %v", err, string(data))
	}

	if !reflect.DeepEqual(data, want) {
		t.Errorf("got %v, expected %v", string(data), want)
	}
}

// Create a server
func runServer(address string, stop chan struct{}) {
	circuitBinary := filepath.Join(os.Getenv("GOPATH"), "bin", "circuit")

	cmd := exec.Command(circuitBinary, "start", "-discover", address)
	cmd.Start()

	// Wait for stop signal
	<-stop

	err := cmd.Process.Signal(syscall.SIGINT)
	if err != nil {
		fmt.Errorf("Release failed: %v", err)
	}

	// Stopped, so continue
	stop <- struct{}{}
}

func stopServer(stop chan struct{}) {
	/* Send stop signal */
	stop <- struct{}{}

	/* Wait until stop succeeded */
	<-stop

	fmt.Println("Stopped")
}

func readData(c *circuit.Client) ([]byte, error) {
	readAnchor := getAnchor(c)
	readChannel := readAnchor.Get().(circuit.Chan)

	if readChannel == nil {
		return nil, fmt.Errorf("Read should have been possible.")
	}

	reader, err := readChannel.Recv()

	if err != nil {
		return nil, fmt.Errorf("Failed to create a reader")
	}

	data, err := ioutil.ReadAll(reader)

	// err should be EOF
	if err == nil || err == io.EOF {
		return data, nil
	}

	return data, err
}
