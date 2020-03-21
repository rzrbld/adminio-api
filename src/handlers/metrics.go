package handlers

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var opsProcessed = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "bucket_size_current",
	Help: "bucket size in kbytes",
}, []string{"bucket"})

func RecordMetrics() {
	go func() {
		for {
			du, err := madmClnt.DataUsageInfo()
			if err != nil {
				log.Print("Error while getting bucket size metrics from server")
			} else {
				if len(du.BucketsSizes) != 0 {
					for k, v := range du.BucketsSizes {
						opsProcessed.WithLabelValues(string(k)).Set(float64(v))
					}
				}
			}
			time.Sleep(30 * time.Minute)
		}
	}()
}
