package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestParseOperator(t *testing.T) {
	testData := []struct {
		name     string
		data     string
		expected Operation
		errMsg   string
	}{
		{"add", "+", Addition, ""},
		{"sub", "-", Subtraction, ""},
		{"mul", "*", Multiplication, ""},
		{"div", "/", Division, ""},
		{"err", "_", InvalidOp, "invalid operation '_'"},
	}

	for _, d := range testData {
		t.Run(d.name, func(t *testing.T) {
			result, err := ParseOperation(d.data)
			if result != d.expected {
				t.Errorf("Expected '%v', got '%v'", d.expected, result)
			}

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}
			if errMsg != d.errMsg {
				t.Errorf("Expected error message '%s', got '%s'", d.errMsg, errMsg)
			}
		})
	}
}

func Test_parser(t *testing.T) {
	testData := []struct {
		name     string
		data     []byte
		expected Input
		errMsg   string
	}{
		{
			"empty",
			[]byte{},
			Input{},
			"data does not contain 4 lines",
		},
		{
			"bad op",
			[]byte("badop\n_\n1\n2"),
			Input{Id: "badop"},
			"invalid operation '_'",
		},
		{
			"bad val1",
			[]byte("badval1\n+\nx\n2"),
			Input{Id: "badval1", Op: Addition},
			"strconv.Atoi: parsing \"x\": invalid syntax",
		},
		{
			"bad val2",
			[]byte("badval2\n+\n1\nx"),
			Input{Id: "badval2", Op: Addition, Val1: 1},
			"strconv.Atoi: parsing \"x\": invalid syntax",
		},
		{
			"good",
			[]byte("good\n+\n1\n2"),
			Input{Id: "good", Op: Addition, Val1: 1, Val2: 2},
			"",
		},
	}

	for _, d := range testData {
		t.Run(d.name, func(t *testing.T) {
			result, err := parser(d.data)
			if result != d.expected {
				t.Errorf("Expected %v, got %v", d.expected, result)
			}

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}
			if errMsg != d.errMsg {
				t.Errorf("Expected error message '%s', got '%s'", d.errMsg, errMsg)
			}
		})
	}
}

func Fuzz_parser(f *testing.F) {
	testcases := [][]byte{
		[]byte("add\n+\n1\n2"),
		[]byte("sub\n-\n1\n2"),
		[]byte("mul\n*\n1\n2"),
		[]byte("div\n/\n1\n2"),
		[]byte("dbz\n/\n1\n0"),
	}

	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, in []byte) {
		_, err := parser(in)
		if err != nil &&
			err.Error() != "data does not contain 4 lines" &&
			!strings.HasPrefix(err.Error(), "invalid operation '") &&
			!errors.Is(err, strconv.ErrSyntax) &&
			!errors.Is(err, strconv.ErrRange) {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestDataParser(t *testing.T) {
	testData := []struct {
		name           string
		data           []byte
		expectedId     string
		expectedValue  int
		expectedErrMsg string
	}{
		{"err", []byte{}, "", 0, "data does not contain 4 lines"},
		{"add", []byte("add\n+\n1\n2"), "add", 3, ""},
		{"sub", []byte("sub\n-\n1\n2"), "sub", -1, ""},
		{"mul", []byte("mul\n*\n1\n2"), "mul", 2, ""},
		{"div", []byte("div\n/\n1\n2"), "div", 0, ""},
		{"dbz", []byte("dbz\n/\n1\n0"), "dbz", 0, "division by zero"},
	}

	for _, d := range testData {
		t.Run(d.name, func(t *testing.T) {
			in := make(chan []byte)
			out := make(chan Result)
			go DataProcessor(in, out)

			in <- d.data
			close(in)

			result := <-out

			if result.Id != d.expectedId {
				t.Errorf("Expected id '%s', got '%s'", d.expectedId, result.Id)
			}

			if result.Value != d.expectedValue {
				t.Errorf("Expected value '%d', got '%d'", d.expectedValue, result.Value)
			}

			var errMsg string
			if result.Err != nil {
				errMsg = result.Err.Error()
			}
			if errMsg != d.expectedErrMsg {
				t.Errorf("Expected error message '%s', got '%s'", d.expectedErrMsg, errMsg)
			}

		})
	}
}

func TestWriteData(t *testing.T) {
	testData := []struct {
		name     string
		data     Result
		expected string
	}{
		{"empty", Result{}, ":0\n"},
		{"nonempty", Result{Id: "x", Value: 1}, "x:1\n"},
	}

	for _, d := range testData {
		t.Run(d.name, func(t *testing.T) {
			in := make(chan Result)
			sb := new(strings.Builder)
			mu := new(sync.Mutex)
			go WriteData(in, sb, mu)

			in <- d.data
			close(in)

			mu.Lock()
			actual := sb.String()
			mu.Unlock()

			if actual != d.expected {
				t.Errorf("Expected string '%s', got '%s'", d.expected, actual)
			}
		})
	}
}

type BadBody struct{}

func (bb BadBody) Read(p []byte) (int, error) {
	return 0, errors.New("intentional bad body")
}

func TestNewController(t *testing.T) {
	testData := []struct {
		name         string
		doBadData    bool
		doBusy       bool
		expectedData []byte
		expectedCode int
		expectedBody string
	}{
		{"empty", false, false, []byte{}, http.StatusAccepted, "OK: 1"},
		{"nonempty", false, false, []byte("nonempty"), http.StatusAccepted, "OK: 1"},
		{"baddata", true, false, []byte("baddata"), http.StatusBadRequest, "Bad Input"},
		{"busy", false, true, []byte("busy"), http.StatusServiceUnavailable, "Too Busy: 1"},
	}

	for _, d := range testData {
		t.Run(d.name, func(t *testing.T) {
			out := make(chan []byte, 1)
			defer close(out)
			handler := NewController(out)

			if d.doBusy {
				out <- d.expectedData
			}

			// Request and Response.
			var body io.Reader = bytes.NewReader(d.expectedData)
			if d.doBadData {
				body = BadBody{}
			}

			req := httptest.NewRequest(http.MethodGet, "/", body)
			rsp := httptest.NewRecorder()
			handler.ServeHTTP(rsp, req)

			// Check tests.
			if rsp.Code != d.expectedCode {
				t.Errorf("Expected status code %d, got %d", d.expectedCode, rsp.Code)
			}

			if body := rsp.Body.String(); body != d.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", d.expectedBody, body)
			}

			if !d.doBadData {
				actualData := <-out
				if slices.Compare(actualData, d.expectedData) != 0 {
					t.Errorf("Expected data '%s', got '%s'", d.expectedData, string(actualData))
				}
			}
		})
	}
}

func Test_run(t *testing.T) {
	testData := []struct {
		name          string
		config        runConfig
		expectBadFile bool
	}{
		{"simple", runConfig{100, 100, t.TempDir() + "/output.txt", "127.0.0.1:8675"}, false},
		{"privFile", runConfig{100, 100, "/output.txt", "127.0.0.1:8675"}, true},
	}

	for _, d := range testData {
		t.Run(d.name, func(t *testing.T) {
			server, chan_err := run(d.config)

			if !d.expectBadFile {
				err := server.Close()
				if err != nil {
					t.Errorf("Error while closing server: %v", err)
				}
			}

			if err := <-chan_err; !errors.Is(err, http.ErrServerClosed) {
				if !d.expectBadFile || !errors.Is(err, os.ErrPermission) {
					t.Errorf("Error from server: %v", err)
				}
			}

		})
	}
}
