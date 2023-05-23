package main

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/version"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":2121").String()
	logLevel      = kingpin.Flag("log.level", "LogLevel - Debug, Info, Warn, Error").Default("Debug").String()

	// Metrics about the exporter itself.
	netDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "fohhnnet_collection_duration_seconds",
			Help: "Duration of collections by the FohhnNet exporter",
		},
		[]string{"module"},
	)
	netRequestErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "fohhnnet_request_errors_total",
			Help: "Errors in requests to the fohhnnet exporter",
		},
	)

	reloadCh chan chan error
)

func init() {
	prometheus.MustRegister(netDuration)
	prometheus.MustRegister(netRequestErrors)
	prometheus.MustRegister(version.NewCollector("fohhnnet_exporter"))
}

func handler(w http.ResponseWriter, r *http.Request, logger log.Logger) {
	query := r.URL.Query()

	target := query.Get("target")
	if len(query["target"]) != 1 || target == "" {
		http.Error(w, "'target' parameter must be specified once", 400)

		return
	}

	logger = log.With(logger, "target", target)
	level.Debug(logger).Log("msg", "Starting scrape")

	start := time.Now()
	registry := prometheus.NewRegistry()
	collector := collector{ctx: r.Context(), target: target, logger: logger}
	registry.MustRegister(collector)
	// Delegate http serving to Prometheus client library, which will call collector.Collect.
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
	duration := time.Since(start).Seconds()
	netDuration.WithLabelValues("FohhnNet").Observe(duration)
	level.Debug(logger).Log("msg", "Finished scrape", "duration_seconds", duration)
}

func main() {

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	level.Info(logger).Log("msg", "Starting fohhnnet_exporter...")

	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	switch lol := *logLevel; lol {
	case "Debug":
		logger = level.NewFilter(logger, level.AllowDebug())
		level.Info(logger).Log("msg", "Starting with loglevel Debug")
	case "Info":
		logger = level.NewFilter(logger, level.AllowInfo())
		level.Info(logger).Log("msg", "Starting with loglevel Info")
	case "Warn":
		logger = level.NewFilter(logger, level.AllowWarn())
		level.Info(logger).Log("msg", "Starting with loglevel Warn")
	case "Error":
		logger = level.NewFilter(logger, level.AllowError())
		level.Info(logger).Log("msg", "Starting with loglevel Error")
	}

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/fohhnnet", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, logger)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
            <head>
            <title>Fohhn-Net Exporter</title>
            <style>
            label{
            display:inline-block;
            width:75px;
            }
            form label {
            margin: 10px;
            }
            form input {
            margin: 10px;
            }
            </style>
            </head>
            <body>
            <h1>Fohhn-Net Exporter</h1>
			<h2>query device</h2>
            <form action="/fohhnnet">
            <label>Target:</label> <input type="text" name="target" placeholder="X.X.X.X" value="1.2.3.4"><br>
            <input type="submit" value="Query">
            </form>
            </body>
            </html>`))
	})

	http.ListenAndServe(*listenAddress, nil)
}
