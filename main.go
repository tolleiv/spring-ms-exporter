package main

import (
	"flag"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strings"
)

var (
	addr = flag.String("listen-address", ":9117", "The address to listen on for HTTP requests.")
)

type Microservice struct {
	Build MicroserviceBuild `json:"build"`
}
type MicroserviceBuild struct {
	Version string `json:"version"`
}
type MicroserviceHealth struct {
	Status string `json:"status"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
            <head><title>Microservice Exporter</title></head>
            <body>
            <h1>Microservice Exporter</h1>
            <p><a href="/probe">Run a probe</a></p>
            <p><a href="/metrics">Metrics</a></p>
            </body>
            </html>`))
	})
	flag.Parse()
	http.HandleFunc("/probe", probeHandler)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func probeHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	target := params.Get("target")
	if target == "" {
		http.Error(w, "Target parameter is missing", 400)
		return
	}
	healthGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "microservice_health",
		Help: "Displays whether or not the service is healthy",
	}, []string{"version"})

	registry := prometheus.NewRegistry()
	registry.MustRegister(healthGauge)

	bytes, err := getJson(target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ms_meta := Microservice{Build: MicroserviceBuild{Version: "unknown"}}
	err = json.Unmarshal([]byte(bytes), &ms_meta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err = getJson(strings.Replace(target, "/info", "/health", 1))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ms_health := MicroserviceHealth{}
	err = json.Unmarshal([]byte(bytes), &ms_health)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	health := 1.0
	if strings.ToUpper(ms_health.Status) != "UP" {
		health = 0
	}
	healthGauge.WithLabelValues(ms_meta.Build.Version).Set(health)

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func getJson(target string) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Get(target)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
