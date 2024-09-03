package client

import (
	"bytes"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SpaceSlow/execenv/internal/config"
	"github.com/SpaceSlow/execenv/internal/metrics"
	"github.com/SpaceSlow/execenv/internal/utils"
)

type httpStrategy struct {
	url  string
	cert *rsa.PublicKey
}

func newHttpStrategy() (*httpStrategy, error) {
	cfg, err := config.GetAgentConfig()
	if err != nil {
		return nil, err
	}
	s := &httpStrategy{
		url: "http://" + cfg.ServerAddr.String() + "/updates/",
	}

	if cfg.CertFile != "" {
		s.cert, err = utils.GetPublicKey(cfg.CertFile)
		if err != nil {
			return nil, fmt.Errorf("extract public key from file error: %w", err)
		}
	}

	return s, nil
}

func (s *httpStrategy) Send(metrics []metrics.Metric) error {
	cfg, err := config.GetAgentConfig()
	if err != nil {
		return err
	}
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	var hash string
	if cfg.Key != "" {
		h := sha256.New()
		h.Write(append(data, []byte(cfg.Key)...))
		hash = hex.EncodeToString(h.Sum(nil))
	}

	data, err = utils.Compress(data)
	if err != nil {
		return err
	}

	if s.cert != nil {
		data, err = rsa.EncryptPKCS1v15(crand.Reader, s.cert, data)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(http.MethodPost, s.url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Real-IP", cfg.LocalIP)

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

	return <-utils.RetryFunc(sendMetrics, cfg.Delays)
}
