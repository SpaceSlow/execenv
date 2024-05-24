package metrics

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
)

type MetricSender struct {
	metrics []Metric
	url     string
	key     string
}

func (ms *MetricSender) Push(metrics []Metric) {
	ms.metrics = make([]Metric, len(metrics))
	copy(ms.metrics, metrics)
}

func (ms *MetricSender) Flush() {
	ms.metrics = nil
}

func (ms *MetricSender) Send() error {
	jsonMetric, err := json.Marshal(ms.metrics)
	if err != nil {
		return err
	}

	var hash string
	if ms.key != "" {
		h := sha256.New()
		h.Write(jsonMetric)
		hash = hex.EncodeToString(h.Sum([]byte(ms.key)))
	}

	req, err := newCompressedRequest(http.MethodPost, ms.url, jsonMetric)
	if err != nil {
		return err
	}

	if hash != "" {
		req.Header.Set("Hash", hash)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.Body.Close() != nil {
		return err
	}

	ms.Flush()
	return nil
}

func NewMetricSender(url, key string) *MetricSender {
	return &MetricSender{
		url: url,
		key: key,
	}
}
