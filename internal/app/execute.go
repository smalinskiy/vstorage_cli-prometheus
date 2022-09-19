package app

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	domain "test.local/prometheus-cli-exporter/internal/repository"
)

type commands struct {
	app, arg0, arg1, arg2, arg3 string
}

type Stats struct {
	Name                 string `xml:"cluster_name"`
	Status               string `xml:"status"`
	License_Status       int    `xml:"license>status"`
	Licensed_Capacity    string `xml:"license>capacity"`
	SpaceAllocatable     int    `xml:"space>allocatable"`
	Space_EffectiveTotal int    `xml:"space>effective_totalclear"`
	SpaceTotal           int    `xml:"space>total"`
	SpaceFree            int    `xml:"space>free"`
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_ops_total",
		Help: "The total number of processed events",
	})
)

var (
	license_status = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "license_status",
		Help: " 0 = NO license installed",
	})
)

var (
	sample_value = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_s3_volume",
		Help: "The total volume of S3 Storage",
	})
)

// func scheduler(wg *sync.WaitGroup, instance int) {
func scheduler(ch chan bool) {
	ticker := time.NewTicker(5 * time.Second)
	fmt.Println("scheduler start")
	for {
		select {
		case <-ticker.C:
			if err := collectstats(); err != nil {
				fmt.Println("Zopa ne rabot. err:", err.Error())
			}
		case <-ch:
			fmt.Println("scheduler stop")
			ticker.Stop()
		}
	}
}

func collectstats() error {
	x := commands{
		app:  "vstorage",
		arg0: "-c",
		arg1: "test1",
		arg2: "stat",
		arg3: "--xml",
	}
	cmd := exec.Command(x.app, x.arg0, x.arg1, x.arg2, x.arg3)
	stdout, _ := cmd.Output()
	var s Stats
	if err := xml.Unmarshal(stdout, &s); err != nil {
		return err
	}
	fmt.Println(s)
	//opsProcessed.Inc()
	sample_value.Set(float64(s.SpaceTotal))
	license_status.Set(float64(s.License_Status))
	return nil
}

func Run(c domain.Config) {
	fmt.Printf("%s", c)
	chStop := make(chan bool)
	go scheduler(chStop)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(c.IP+":"+c.Port, nil)
	if err != nil {
		log.Fatalf("http.ListenAndServer: %v\n", err)
	}
	chStop <- true

}
