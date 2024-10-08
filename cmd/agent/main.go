package main

import (
	"github.com/SpaceSlow/execenv/internal/worker"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SpaceSlow/execenv/internal/config"
	"github.com/SpaceSlow/execenv/internal/metrics"
)

func main() {
	config.PrintBuildInfo()
	cfg, err := config.GetAgentConfig()
	if err != nil {
		log.Fatalf("stopped agent: %s", err)
	}

	pollTick := time.Tick(cfg.PollInterval.Duration)
	reportTick := time.Tick(cfg.ReportInterval.Duration)
	metricWorkers, err := worker.NewMetricWorkers()
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
			sendCh <- <-pollCh
			if err = <-metricWorkers.Err(); err != nil {
				log.Println(err)
			} else {
				log.Println("sended metrics")
			}
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
		case err = <-metricWorkers.Err():
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println("sended metrics")
		}
	}
}
