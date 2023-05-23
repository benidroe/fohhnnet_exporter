package main

import (
	"context"
	"fohhnnet_exporter/FohhnNet"
	"github.com/go-kit/kit/log"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type collector struct {
	ctx    context.Context
	target string
	//module *config.Module
	logger log.Logger
}

// Describe implements Prometheus.Collector.
func (c collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("dummy", "dummy", nil, nil)
}

// Collect implements Prometheus.Collector.
func (c collector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	var pjSlice []prometheus.Metric // place to push collected metrics

	walkFohhnNet(c.target, &pjSlice, c.logger)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc("fohhnnet_walk_duration_seconds", "Time FohhnNet walk took.", nil, nil),
		prometheus.GaugeValue,
		time.Since(start).Seconds())

	// iterate over results and push results to chan
	for i := range pjSlice {
		ch <- pjSlice[i]
	}

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc("fohhnnet_scrape_duration_seconds", "Total FohhnNet time scrape took (walk and processing).", nil, nil),
		prometheus.GaugeValue,
		time.Since(start).Seconds())
}

func walkFohhnNet(dest string, pjSlice *[]prometheus.Metric, logger log.Logger) {

	fohhnNetSession, err := FohhnNet.NewFohhnNetTcpSession(dest, 4001)
	hasAnswered := false
	_ = hasAnswered

	isUp := 0

	if err == nil {

		defer FohhnNet.Close(fohhnNetSession) // finally, close the connection

		if fohhnNetSession.IsConnected {

			isUp = 1
			res := FohhnNet.ScanFohhnNet(fohhnNetSession, 1, 5)

			for _, id := range res {

				walkResult, err := FohhnNet.ScrapeFohhnDevice(fohhnNetSession, id)
				if err == nil {

					hasAnswered = true

					// Append Metric to result set pjSlice
					*pjSlice = append(*pjSlice, prometheus.MustNewConstMetric(
						prometheus.NewDesc("fohhnnet_up", "fohhnNet device is up", []string{"id", "model", "version"}, nil),
						prometheus.GaugeValue,
						float64(1), strconv.Itoa(int(id)), FohhnNet.GetModelNameByNumber(walkResult.Device), walkResult.Version))

					// Append Metric to result set pjSlice
					*pjSlice = append(*pjSlice, prometheus.MustNewConstMetric(
						prometheus.NewDesc("fohhnnet_temperature", "fohhnNet device temperature", []string{"id"}, nil),
						prometheus.GaugeValue,
						float64(int(walkResult.Temperature*10)), strconv.Itoa(int(id))))

					if len(walkResult.Protect) == len(walkResult.OutputChannelName) && len(walkResult.Protect) == len(walkResult.SpeakerPreset) {
						for k, v := range walkResult.Protect {
							protect := 1
							if !v {
								protect = 0
							}
							// Append Metric to result set pjSlice
							*pjSlice = append(*pjSlice, prometheus.MustNewConstMetric(
								prometheus.NewDesc("fohhnnet_protect", "fohhnNet device channel protection", []string{"id", "channel", "name", "preset"}, nil),
								prometheus.GaugeValue,
								float64(protect), strconv.Itoa(int(id)), strconv.Itoa(int(k+1)), strings.TrimSpace(walkResult.OutputChannelName[k]), strings.TrimSpace(walkResult.SpeakerPreset[k])))
						}
					}

					// Append Metric to result set pjSlice
					*pjSlice = append(*pjSlice, prometheus.MustNewConstMetric(
						prometheus.NewDesc("fohhnnet_operatinghours", "fohhnNet device operating time in hours", []string{"id"}, nil),
						prometheus.GaugeValue,
						float64(int(walkResult.OperatingTimeHours)), strconv.Itoa(int(id))))

					power := 1
					if walkResult.Standby {
						power = 0
					}
					// Append Metric to result set pjSlice
					*pjSlice = append(*pjSlice, prometheus.MustNewConstMetric(
						prometheus.NewDesc("fohhnnet_power", "fohhnNet device standby power", []string{"id"}, nil),
						prometheus.GaugeValue,
						float64(power), strconv.Itoa(int(id))))

				}

			}
		}

	}

	// Append Metric to result set pjSlice
	*pjSlice = append(*pjSlice, prometheus.MustNewConstMetric(
		prometheus.NewDesc("fohhnnet_adapter_up", "fohhnNet adapter is up", nil, nil),
		prometheus.GaugeValue,
		float64(isUp)))

}
