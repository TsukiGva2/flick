package flick

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/MyTempoESP/serial"
)

type Forth struct {
	port         *serial.Port
	mu           sync.Mutex
	responseChan chan string
}

func NewForth(dev string, timeout time.Duration) (f Forth, err error) {

	conf := &serial.Config{
		Name:        dev,
		Baud:        115200,
		ReadTimeout: timeout,
	}

	f.port, err = serial.OpenPort(conf)

	if err != nil {

		log.Fatalf("Failed to open serial port: %v", err)
	}

	f.responseChan = make(chan string)

	return
}

func (f *Forth) Stop() {

	f.port.Close()
	close(f.responseChan)
}

func (f *Forth) Start() {

	go func() {

		buf := make([]byte, 128)

		for {
			n, err := f.port.Read(buf)

			if err != nil {

				f.responseChan <- "(timeout!)"

				continue
			}

			if n > 0 {

				f.responseChan <- string(buf[:n])
			}
		}
	}()
}

func (f *Forth) Send(input string) (response string, err error) {

	f.mu.Lock()
	defer f.mu.Unlock()

	_, err = f.port.Write([]byte(input + "\n"))

	if err != nil {

		log.Printf("Failed to send data: %v", err)

		return
	}

	response = <-f.responseChan

	return
}

func (f *Forth) Query(input string) (multilineResponse string, err error) {

	f.mu.Lock()
	defer f.mu.Unlock()

	_, err = f.port.Write([]byte(input + "\n"))

	if err != nil {

		log.Printf("Failed to send data: %v", err)

		return
	}

	fmt.Printf("Sent: %s\n", input)

	for {
		response := <-f.responseChan
		multilineResponse += response

		if strings.Contains(multilineResponse, "ok") {

			break
		}
	}

	return
}
