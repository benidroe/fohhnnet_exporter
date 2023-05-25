package main

import (
	"context"
	"fmt"
	"fohhnnet_exporter/FohhnNet"
	"github.com/go-kit/kit/log"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type collector struct {
	ctx      context.Context
	target   string
	port     int
	protocol string
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

	walkFohhnNet(c.target, c.port, c.protocol, &pjSlice, c.logger)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc("fohhnnet_walk_duration_seconds", "Time Fohhn-Net walk took.", nil, nil),
		prometheus.GaugeValue,
		time.Since(start).Seconds())

	// iterate over results and push results to chan
	for i := range pjSlice {
		ch <- pjSlice[i]
	}

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc("fohhnnet_scrape_duration_seconds", "Total Fohhn-Net time scrape took (walk and processing).", nil, nil),
		prometheus.GaugeValue,
		time.Since(start).Seconds())
}

func walkFohhnNet(dest string, port int, protocol string, pjSlice *[]prometheus.Metric, logger log.Logger) {

	fohhnNetSession, err := FohhnNet.NewFohhnNetSession(dest, port, protocol)
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
						prometheus.NewDesc("fohhnnet_up", "Fohhn-Net device is up", []string{"id", "model", "version"}, nil),
						prometheus.GaugeValue,
						float64(1), strconv.Itoa(int(id)), FohhnNet.GetModelNameByNumber(walkResult.Device), walkResult.Version))

					// Append Metric to result set pjSlice
					*pjSlice = append(*pjSlice, prometheus.MustNewConstMetric(
						prometheus.NewDesc("fohhnnet_temperature", "Fohhn-Net device temperature", []string{"id"}, nil),
						prometheus.GaugeValue,
						float64(int(walkResult.Temperature)), strconv.Itoa(int(id))))

					fmt.Println("ABC", len(walkResult.Protect), len(walkResult.OutputChannelName), len(walkResult.SpeakerPreset))

					if len(walkResult.Protect) >= walkResult.NumOfChannels && len(walkResult.OutputChannelName) >= walkResult.NumOfChannels && len(walkResult.SpeakerPreset) >= walkResult.NumOfChannels {
						for k, _ := range walkResult.OutputChannelName {
							protect := 1
							if walkResult.Protect[k] {
								protect = 0
							}
							// Append Metric to result set pjSlice
							*pjSlice = append(*pjSlice, prometheus.MustNewConstMetric(
								prometheus.NewDesc("fohhnnet_protect", "Fohhn-Net device channel protection", []string{"id", "channel", "name", "preset"}, nil),
								prometheus.GaugeValue,
								float64(protect), strconv.Itoa(int(id)), strconv.Itoa(int(k+1)), strings.TrimSpace(walkResult.OutputChannelName[k]), strings.TrimSpace(walkResult.SpeakerPreset[k])))
						}
					}

					// Append Metric to result set pjSlice
					*pjSlice = append(*pjSlice, prometheus.MustNewConstMetric(
						prometheus.NewDesc("fohhnnet_operatinghours", "Fohhn-Net device operating time in hours", []string{"id"}, nil),
						prometheus.GaugeValue,
						float64(int(walkResult.OperatingTimeHours)), strconv.Itoa(int(id))))

					power := 1
					if walkResult.Standby {
						power = 0
					}
					// Append Metric to result set pjSlice
					*pjSlice = append(*pjSlice, prometheus.MustNewConstMetric(
						prometheus.NewDesc("fohhnnet_power", "Fohhn-Net device standby power", []string{"id"}, nil),
						prometheus.GaugeValue,
						float64(power), strconv.Itoa(int(id))))

				}

			}
		}

	}

	// Append Metric to result set pjSlice
	*pjSlice = append(*pjSlice, prometheus.MustNewConstMetric(
		prometheus.NewDesc("fohhnnet_adapter_up", "Fohhn-Net adapter is up", nil, nil),
		prometheus.GaugeValue,
		float64(isUp)))

}
