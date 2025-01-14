package flick

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/MyTempoESP/serial"
)

// LABELS
const (
	PORTAL = iota
	UNICAS
	REGIST
	COMUNICANDO
	LEITOR
	LTE4G
	WIFI
	IP
	LOCAL
	PROVA

	LABELS_COUNT
)

// VALUES
const (
	WEB = iota
	CONECTAD
	DESLIGAD
	AUTOMATIC
	OK
	X

	VALUES_COUNT
)

type MyTempo_Forth struct {
	port         *serial.Port
	mu           sync.Mutex
	responseChan chan string
}

func NewForth(dev string) (f MyTempo_Forth, err error) {

	// Configure the serial port
	conf := &serial.Config{
		Name:        dev,    // Update to match your serial port
		Baud:        115200, // Adjust the baud rate as needed
		ReadTimeout: time.Second * 1,
	}

	// Open the serial port
	f.port, err = serial.OpenPort(conf)

	if err != nil {

		log.Fatalf("Failed to open serial port: %v", err)
	}

	f.responseChan = make(chan string)

	return
}

func (f *MyTempo_Forth) Stop() {

	f.port.Close()
	close(f.responseChan)
}

func (f *MyTempo_Forth) Start() {

	// Goroutine to read data from the Arduino
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

func (f *MyTempo_Forth) Send(input string) (response string, err error) {

	f.mu.Lock()
	defer f.mu.Unlock()

	_, err = f.port.Write([]byte(input + "\n"))

	if err != nil {

		log.Printf("Failed to send data: %v", err)

		return
	}

	fmt.Printf("Sent: %s\n", input)

	// Wait for a response with synchronization
	response = <-f.responseChan

	fmt.Printf("> %s\n", strings.TrimSpace(response))

	//time.Sleep(50 * time.Millisecond)

	return
}

func (f *MyTempo_Forth) Consume() (output string) {

	for {

		select {
		case resp := <-f.responseChan:
			output += resp
		case <-time.After(2 * time.Second):
			break
		}
	}

	return
}
