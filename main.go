package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var prefix = "sample_external_url"

var (
	reg = prometheus.NewRegistry()

	status = promauto.With(reg).NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prefix + "_up",
			Help: "Boolean status of site up or down",
		}, []string{"url"},
	)
	latency = promauto.With(reg).NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prefix + "_response_ms",
			Help: "HTTP response in milliseconds",
		}, []string{"url"},
	)
)

var netClient = &http.Client{
	Timeout: time.Second * 5,
}

// makeRequest - sets a Prometheus metric for response time in ms and status code
func makeRequest(url string) error {
	start := time.Now()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := netClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	duration := time.Since(start).Milliseconds()

	upOrDown, err := GetStatus(resp.StatusCode)
	if err != nil {
		return err
	}

	status.With(prometheus.Labels{"url": url}).Set(float64(upOrDown))
	latency.With(prometheus.Labels{"url": url}).Set(float64(duration))

	return nil
}

func main() {

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	config, err := LoadConfig("conf.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, uri := range config.URLs {
		go func(u string) {
			for {
				time.Sleep(time.Second * 5)
				makeRequest(u)
			}
		}(uri)
	}

	log.Fatal(http.ListenAndServe(":8080", nil))

}
