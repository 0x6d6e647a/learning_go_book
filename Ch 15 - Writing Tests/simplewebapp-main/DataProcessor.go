package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Operation uint8

const (
	InvalidOp Operation = iota
	Addition
	Subtraction
	Multiplication
	Division
)

func ParseOperation(s string) (Operation, error) {
	switch s {
	case "+":
		return Addition, nil
	case "-":
		return Subtraction, nil
	case "*":
		return Multiplication, nil
	case "/":
		return Division, nil
	}

	return InvalidOp, fmt.Errorf("invalid operation '%s'", s)
}

type Input struct {
	Id   string
	Op   Operation
	Val1 int
	Val2 int
}

func parser(data []byte) (input Input, err error) {
	lines := bytes.Split(data, []byte("\n"))
	if len(lines) != 4 {
		return input, errors.New("data does not contain 4 lines")
	}

	input.Id = string(lines[0])

	input.Op, err = ParseOperation(string(lines[1]))
	if err != nil {
		return input, err
	}

	input.Val1, err = strconv.Atoi(string(lines[2]))
	if err != nil {
		return input, err
	}

	input.Val2, err = strconv.Atoi(string(lines[3]))
	if err != nil {
		return input, err
	}

	return input, nil
}

type Result struct {
	Id    string
	Value int
	Err   error
}

func DataProcessor(in <-chan []byte, out chan<- Result) {
	for data := range in {
		var result Result

		input, err := parser(data)
		if err != nil {
			result.Err = err
			out <- result
			continue
		}

		result.Id = input.Id

		switch input.Op {
		case Addition:
			result.Value = input.Val1 + input.Val2
		case Subtraction:
			result.Value = input.Val1 - input.Val2
		case Multiplication:
			result.Value = input.Val1 * input.Val2
		case Division:
			if input.Val2 == 0 {
				result.Err = errors.New("division by zero")
				out <- result
				continue
			}
			result.Value = input.Val1 / input.Val2
		}

		out <- result
	}

	close(out)
}

func WriteData(in <-chan Result, w io.Writer, mu *sync.Mutex) {
	for r := range in {
		mu.Lock()
		fmt.Fprintf(w, "%s:%d\n", r.Id, r.Value)
		mu.Unlock()
	}
}

func NewController(out chan<- []byte) http.Handler {
	var numRecieved int
	var numRejected int

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		numRecieved += 1

		data, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			numRejected += 1
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Input"))
			return
		}

		// Write it to the queue in raw format.
		select {
		case out <- data:
			// Success!
		default:
			// If the channel is backed up, return an error.
			numRejected += 1
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Too Busy: " + strconv.Itoa(numRejected)))
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("OK: " + strconv.Itoa(numRecieved)))
	})
}

type runConfig struct {
	byteChanLen   int
	resultChanLen int
	outputFile    string
	addr          string
}

func run(rc runConfig) (*http.Server, <-chan error) {
	chan_err := make(chan error, 1)

	// Data processor.
	chan_byte := make(chan []byte, rc.byteChanLen)
	chan_result := make(chan Result, rc.resultChanLen)
	go DataProcessor(chan_byte, chan_result)

	// Output file.
	f, err := os.Create(rc.outputFile)
	if err != nil {
		chan_err <- err
		return nil, chan_err
	}

	mu := new(sync.Mutex)
	go WriteData(chan_result, f, mu)

	// HTTP server.
	server := &http.Server{
		Addr:    rc.addr,
		Handler: NewController(chan_byte),
	}

	go func() {
		defer close(chan_err)
		defer close(chan_byte)
		defer f.Close()

		err = server.ListenAndServe()
		chan_err <- err
	}()

	return server, chan_err
}

func main() {
	mainConfig := runConfig{
		100, 100, "results.txt", ":8080",
	}

	_, chan_err := run(mainConfig)

	for err := range chan_err {
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			break
		}
		panic(err)
	}
}
