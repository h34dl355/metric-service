package internal

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type gauge float64
type counter int64

var pollCount int

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetStats(t time.Duration, metrics *map[string]interface{}) {
	ticker := time.NewTicker(t)
	for {
		select {
		case <-ticker.C:
			pollCount++
			collectStats(metrics)
		default:
			continue
		}
	}
}

func collectStats(metrics *map[string]interface{}) {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	*metrics = map[string]interface{}{
		"Alloc":         gauge(m.Alloc),
		"BuckHashSys":   gauge(m.BuckHashSys),
		"Frees":         gauge(m.Frees),
		"GCCPUFraction": gauge(m.GCCPUFraction),
		"GCSys":         gauge(m.GCSys),
		"HeapAlloc":     gauge(m.HeapAlloc),
		"HeapIdle":      gauge(m.HeapIdle),
		"HeapInuse":     gauge(m.HeapInuse),
		"HeapObjects":   gauge(m.HeapObjects),
		"HeapReleased":  gauge(m.HeapReleased),
		"HeapSys":       gauge(m.HeapSys),
		"LastGC":        gauge(m.LastGC),
		"Lookups":       gauge(m.Lookups),
		"MCacheInuse":   gauge(m.MCacheInuse),
		"MCacheSys":     gauge(m.MCacheSys),
		"MSpanInuse":    gauge(m.MSpanInuse),
		"MSpanSys":      gauge(m.MSpanSys),
		"Mallocs":       gauge(m.Mallocs),
		"NextGC":        gauge(m.NextGC),
		"NumForcedGC":   gauge(m.NumForcedGC),
		"NumGC":         gauge(m.NumGC),
		"OtherSys":      gauge(m.OtherSys),
		"PauseTotalNs":  gauge(m.PauseTotalNs),
		"StackInuse":    gauge(m.StackInuse),
		"StackSys":      gauge(m.StackSys),
		"Sys":           gauge(m.Sys),
		"TotalAlloc":    gauge(m.TotalAlloc),
		"PollCount":     counter(pollCount),
		"RandomValue":   gauge(rand.Float64()),
	}
	fmt.Println((*metrics)["RandomValue"])
}

func SendResponse(metrics *map[string]interface{}, host string) {
	for i, v := range *metrics {
		valueType := "gauge"
		if i == "PollCount" {
			valueType = "counter"
		}
		postRequest := fmt.Sprintf("%s/%s/%s/%v", host, valueType, i, v)
		//fmt.Println(postRequest)
		request, err := http.NewRequest("MethodPost", postRequest, nil)
		if err != nil {
			log.Fatal(err)
		}
		request.Header.Set("Content-Type", "text/plain")
		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
	}
}
