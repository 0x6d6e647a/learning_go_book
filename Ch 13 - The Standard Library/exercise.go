package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

var logger = slog.Default()

type TimeHandlerMultiplexer struct {
	*http.ServeMux
}

func NewTimeHandlerMultiplexer() *TimeHandlerMultiplexer {
	return &TimeHandlerMultiplexer{http.NewServeMux()}
}

func (mux *TimeHandlerMultiplexer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger.Info("request", "addr", req.RemoteAddr, "method", req.Method, "path", req.URL.RawPath)
	mux.ServeMux.ServeHTTP(rw, req)
}

type TimeJSON struct {
	DayOfWeek  string `json:"day_of_week"`
	DayOfMonth int    `json:"day_of_month"`
	Month      string `json:"month"`
	Year       int    `json:"year"`
	Hour       int    `json:"hour"`
	Minute     int    `json:"minute"`
	Second     int    `json:"second"`
}

func NewTimeJSON(t time.Time) TimeJSON {
	return TimeJSON{
		DayOfWeek:  t.Weekday().String(),
		DayOfMonth: t.Day(),
		Month:      t.Month().String(),
		Year:       t.Year(),
		Hour:       t.Hour(),
		Minute:     t.Minute(),
		Second:     t.Second(),
	}
}

func main() {
	// Setup multiplexer.
	handler := NewTimeHandlerMultiplexer()

	handler.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		// Only accept GET method.
		if req.Method != http.MethodGet {
			logger.Warn("invalid method")
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Format time for output.
		now := time.Now()
		var out []byte

		if req.Header.Get("Accept") == "application/json" {
			json, err := json.Marshal(NewTimeJSON(now))
			if err != nil {
				logger.Error("json marshaling error", "err", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			out = json
		} else {
			str := now.Format(time.RFC3339)
			out = []byte(str)
		}

		// Write response.
		rw.WriteHeader(http.StatusOK)
		rw.Write(out)
	})

	// Configure HTTP server.
	server := http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  40 * time.Second,
		Handler:      handler,
	}

	// Launch HTTP server.
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
