// Package monitor: about prometheus and friends
// Good Article: https://www.opsdash.com/blog/golang-app-monitoring-statsd-expvar-prometheus.html
package monitor

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func increment(counter prometheus.Counter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		counter.Add(1)
		msg := fmt.Sprintf("New count: %v", counter)
		log.Println(msg)
		if _, err := w.Write([]byte(msg)); err != nil {
			log.Print(err.Error())
		}
	}
}
func RunPrometheus() {
	var fooCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "foo_total",
		Help: "Number of foo successfully processed.",
	})
	prometheus.MustRegister(fooCount)

	http.HandleFunc("/inc", increment(fooCount))
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Check out http://localhost:9090/metrics or /inc endpoints")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
