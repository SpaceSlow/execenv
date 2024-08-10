package metrics

import (
	"bytes"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// MetricWorkers служит для аккумуляции и отправки метрик на сервер, с заданным ключом.
type MetricWorkers struct {
	metricsForSend chan []Metric
	errorsCh       chan error

	url       string
	key       string
	cert      *rsa.PublicKey
	delays    []time.Duration
	pollCount atomic.Int64
}

func NewMetricWorkers(numWorkers int, url, key, certFile string, delays []time.Duration) (*MetricWorkers, error) {
	var (
		cert *rsa.PublicKey
		err  error
	)
	if certFile != "" {
		cert, err = getPublicKey(certFile)
		if err != nil {
			return nil, fmt.Errorf("extract public key from file error: %w", err)
		}
	}
	return &MetricWorkers{
		metricsForSend: make(chan []Metric, numWorkers),
		errorsCh:       make(chan error, numWorkers),
		url:            url,
		key:            key,
		cert:           cert,
		delays:         delays,
	}, nil
}

func getPublicKey(file string) (*rsa.PublicKey, error) {
	certBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	certBlock, _ := pem.Decode(certBytes)
	return x509.ParsePKCS1PublicKey(certBlock.Bytes)
}

func (mw *MetricWorkers) Send(metrics []Metric) {
	pollCount := mw.pollCount.Load()
	data, err := json.Marshal(metrics)
	if err != nil {
		mw.errorsCh <- err
		return
	}

	var hash string
	if mw.key != "" {
		h := sha256.New()
		h.Write(data)
		hash = hex.EncodeToString(h.Sum([]byte(mw.key)))
	}

	data, err = compress(data)
	if err != nil {
		mw.errorsCh <- err
		return
	}

	if mw.cert != nil {
		data, err = rsa.EncryptPKCS1v15(crand.Reader, mw.cert, data)
		if err != nil {
			mw.errorsCh <- err
			return
		}
	}

	req, err := http.NewRequest(http.MethodPost, mw.url, bytes.NewReader(data))
	if err != nil {
		mw.errorsCh <- err
		return
	}
	setCompressHeader(req)

	if hash != "" {
		req.Header.Set("Hash", hash)
	}

	resCh := make(chan *http.Response, 1)
	defer close(resCh)
	sendMetrics := func() error {
		var res *http.Response
		res, err = http.DefaultClient.Do(req)
		if err != nil {
			if len(resCh) > 0 {
				<-resCh
			}
			resCh <- res
			return err
		}
		if len(resCh) > 0 {
			<-resCh
		}
		if err = res.Body.Close(); err != nil {
			return err
		}
		resCh <- res
		return nil
	}
	err = <-RetryFunc(sendMetrics, mw.delays)

	if err != nil {
		mw.errorsCh <- err
		return
	}
	mw.pollCount.Add(-pollCount)
	mw.errorsCh <- nil
}

func (mw *MetricWorkers) Poll(pollCh chan []Metric) {
	metricSlice := make([]Metric, 0)
	runtimeMetricsCh := mw.getRuntimeMetrics()
	gopsutilMetricsCh := mw.getGopsutilMetrics()
	for m := range fanIn(runtimeMetricsCh, gopsutilMetricsCh) {
		metricSlice = append(metricSlice, m...)
	}
	metricSlice = append(
		metricSlice,
		Metric{
			Type:  Gauge,
			Name:  "RandomValue",
			Value: rand.Float64(),
		},
	)
	metricSlice = append(metricSlice, Metric{
		Type:  Counter,
		Name:  "PollCount",
		Value: mw.pollCount.Add(1),
	})

	if len(pollCh) > 0 {
		<-pollCh
	}
	pollCh <- metricSlice
}

func (mw *MetricWorkers) getGopsutilMetrics() chan []Metric {
	metricsCh := make(chan []Metric)

	go func() {
		defer close(metricsCh)
		v, _ := mem.VirtualMemory()
		cpu, _ := cpu.Percent(0, false)

		metrics := []Metric{
			{
				Type:  Gauge,
				Name:  "TotalMemory",
				Value: float64(v.Total),
			},
			{
				Type:  Gauge,
				Name:  "FreeMemory",
				Value: float64(v.Free),
			},
			{
				Type:  Gauge,
				Name:  "CPUtilization1",
				Value: cpu[0],
			},
		}
		metricsCh <- metrics
	}()

	return metricsCh
}

func (mw *MetricWorkers) getRuntimeMetrics() chan []Metric {
	metricsCh := make(chan []Metric)

	go func() {
		defer close(metricsCh)
		var rtm runtime.MemStats
		runtime.ReadMemStats(&rtm)

		metrics := []Metric{
			{
				Type:  Gauge,
				Name:  "Alloc",
				Value: float64(rtm.Alloc),
			},
			{
				Type:  Gauge,
				Name:  "BuckHashSys",
				Value: float64(rtm.BuckHashSys),
			},
			{
				Type:  Gauge,
				Name:  "Frees",
				Value: float64(rtm.Frees),
			},
			{
				Type:  Gauge,
				Name:  "GCCPUFraction",
				Value: float64(rtm.GCCPUFraction),
			},
			{
				Type:  Gauge,
				Name:  "GCSys",
				Value: float64(rtm.GCSys),
			},
			{
				Type:  Gauge,
				Name:  "HeapAlloc",
				Value: float64(rtm.HeapAlloc),
			},
			{
				Type:  Gauge,
				Name:  "HeapIdle",
				Value: float64(rtm.HeapIdle),
			},
			{
				Type:  Gauge,
				Name:  "HeapInuse",
				Value: float64(rtm.HeapInuse),
			},
			{
				Type:  Gauge,
				Name:  "HeapObjects",
				Value: float64(rtm.HeapObjects),
			},
			{
				Type:  Gauge,
				Name:  "HeapReleased",
				Value: float64(rtm.HeapReleased),
			},
			{
				Type:  Gauge,
				Name:  "HeapSys",
				Value: float64(rtm.HeapSys),
			},
			{
				Type:  Gauge,
				Name:  "LastGC",
				Value: float64(rtm.LastGC),
			},
			{
				Type:  Gauge,
				Name:  "Lookups",
				Value: float64(rtm.Lookups),
			},
			{
				Type:  Gauge,
				Name:  "MCacheInuse",
				Value: float64(rtm.MCacheInuse),
			},
			{
				Type:  Gauge,
				Name:  "MCacheSys",
				Value: float64(rtm.MCacheSys),
			},
			{
				Type:  Gauge,
				Name:  "MSpanInuse",
				Value: float64(rtm.MSpanInuse),
			},
			{
				Type:  Gauge,
				Name:  "MSpanSys",
				Value: float64(rtm.MSpanSys),
			},
			{
				Type:  Gauge,
				Name:  "Mallocs",
				Value: float64(rtm.Mallocs),
			},
			{
				Type:  Gauge,
				Name:  "NextGC",
				Value: float64(rtm.NextGC),
			},
			{
				Type:  Gauge,
				Name:  "NumForcedGC",
				Value: float64(rtm.NumForcedGC),
			},
			{
				Type:  Gauge,
				Name:  "NumGC",
				Value: float64(rtm.NumGC),
			},
			{
				Type:  Gauge,
				Name:  "OtherSys",
				Value: float64(rtm.OtherSys),
			},
			{
				Type:  Gauge,
				Name:  "PauseTotalNs",
				Value: float64(rtm.PauseTotalNs),
			},
			{
				Type:  Gauge,
				Name:  "StackInuse",
				Value: float64(rtm.StackInuse),
			},
			{
				Type:  Gauge,
				Name:  "PauseTotalNs",
				Value: float64(rtm.PauseTotalNs),
			},
			{
				Type:  Gauge,
				Name:  "StackSys",
				Value: float64(rtm.StackSys),
			},
			{
				Type:  Gauge,
				Name:  "Sys",
				Value: float64(rtm.Sys),
			},
			{
				Type:  Gauge,
				Name:  "TotalAlloc",
				Value: float64(rtm.TotalAlloc),
			},
		}
		metricsCh <- metrics
	}()

	return metricsCh
}

func (mw *MetricWorkers) Close() {
	close(mw.errorsCh)
}

func (mw *MetricWorkers) Err() chan error {
	return mw.errorsCh
}
