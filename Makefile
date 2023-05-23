build:
		go get github.com/go-kit/kit/log
		go get github.com/go-kit/kit/log/level
		go get github.com/prometheus/common/version
		go get github.com/prometheus/client_golang/prometheus
		go get github.com/prometheus/client_golang/prometheus/promhttp
		go get github.com/alecthomas/kingpin/v2

		go build fohhnnet_exporter

run:
		go run main.go
