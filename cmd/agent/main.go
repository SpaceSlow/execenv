package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SpaceSlow/execenv/cmd/metrics"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func printBuildInfo() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)
}

func main() {
	printBuildInfo()
	cfg, err := GetConfigWithFlags(os.Args[0], os.Args[1:])

	if err != nil {
		panic(err)
	}

	url := "http://" + cfg.ServerAddr.String() + "/updates/"
	pollInterval := time.Duration(cfg.PollInterval) * time.Second
	reportInterval := time.Duration(cfg.ReportInterval) * time.Second

	pollTick := time.Tick(pollInterval)
	reportTick := time.Tick(reportInterval)
	metricWorkers := metrics.NewMetricWorkers(cfg.RateLimit, url, cfg.Key, cfg.Delays)
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
