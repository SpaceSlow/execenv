package metrics

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func ParseMetricType(mType string) (MetricType, error) {
	switch mType {
	case "counter":
		return Counter, nil
	case "gauge":
		return Gauge, nil
	default:
		return MetricType(-1), ErrIncorrectMetricTypeOrValue
	}
}

func RetryFunc(f func() error, delays []time.Duration) chan error {
	errorCh := make(chan error)

	go func() {
		defer close(errorCh)
		var err error
		for attempt := 0; attempt < len(delays); attempt++ {
			if err = f(); err == nil {
				errorCh <- nil
				return
			}
			<-time.After(delays[attempt])
		}
		errorCh <- err
	}()

	return errorCh
}

func newCompressedRequest(method, url string, data []byte) (*http.Request, error) {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %v", err)
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	req, err := http.NewRequest(method, url, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func fanIn(chs ...chan []Metric) chan []Metric {
	var wg sync.WaitGroup
	outCh := make(chan []Metric)

	output := func(c chan []Metric) {
		for m := range c {
			outCh <- m
		}
		wg.Done()
	}

	wg.Add(len(chs))
	for _, c := range chs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(outCh)
	}()

	return outCh
}
