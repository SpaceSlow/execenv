package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SpaceSlow/execenv/cmd/config"
	"github.com/SpaceSlow/execenv/cmd/metrics"
)

func main() {
	config.PrintBuildInfo()
	cfg, err := config.GetAgentConfig()
	if err != nil {
		log.Fatalf("stopped agent: %s", err)
	}

	url := "http://" + cfg.ServerAddr.String() + "/updates/"
	pollInterval := time.Duration(cfg.PollInterval) * time.Second
	reportInterval := time.Duration(cfg.ReportInterval) * time.Second

	pollTick := time.Tick(pollInterval)
	reportTick := time.Tick(reportInterval)
	metricWorkers, err := metrics.NewMetricWorkers(cfg.RateLimit, url, cfg.Key, cfg.CertFile, cfg.Delays)
	if err != nil {
		log.Fatalf("stopped agent: %s", err)
	}
	sendCh := make(chan []metrics.Metric, cfg.RateLimit)
	pollCh := make(chan []metrics.Metric, 1)

	for w := 0; w < cfg.RateLimit; w++ {
		go func(sendCh chan []metrics.Metric) {
			for ms := range sendCh {
				go metricWorkers.Send(ms)
			}
		}(sendCh)
	}

	closed := make(chan os.Signal, 1)
	defer close(closed)
	signal.Notify(closed, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	for {
		select {
		case <-closed:
			close(sendCh)
			close(pollCh)
			metricWorkers.Close()
			log.Fatal("stopped agent")
		case <-pollTick:
			go metricWorkers.Poll(pollCh)
			log.Println("polled metrics")
		case <-reportTick:
			metricSlice := <-pollCh
			sendCh <- metricSlice
		case err := <-metricWorkers.Err():
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println("sended metrics")
		}
	}
}
