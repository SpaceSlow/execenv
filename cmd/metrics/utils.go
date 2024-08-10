package metrics

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"sync"
	"time"
)

// RetryFunc выполняет функцию f, в случае ошибки последовательно спустя промежутки delays пробует заново,
// в случае неудачи всех попыток кладет в chan error последнюю полученную ошибку, при успехе nil.
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

func compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %w", err)
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %w", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %w", err)
	}
	return b.Bytes(), nil
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
