package utils

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
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

func Compress(data []byte) ([]byte, error) {
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

func OutboundIP(serverAddr string) (string, error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return "", nil
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.TCPAddr)

	return localAddr.IP.String(), nil
}
