package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/h34dl355/metric-service/internal"
)

const (
	pollInterval   time.Duration = 2 * time.Second
	reportInterval time.Duration = 10 * time.Second
	host           string        = "http://127.0.0.1:8080/update"
)

var metrics = make(map[string]interface{})

func main() {
	ctx2, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()

	go internal.GetStats(pollInterval, &metrics)
	ticker := time.NewTicker(reportInterval)
	for {
		select {
		case <-ticker.C:
			internal.SendResponse(&metrics, host)
		case <-ctx2.Done():
			fmt.Println("exit")
			log.Fatal()
		}
	}
}
