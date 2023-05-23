config:
		go get github.com/go-kit/kit/log
		go get github.com/go-kit/kit/log/level
		go get github.com/prometheus/common/version
		go get github.com/prometheus/client_golang/prometheus
		go get github.com/prometheus/client_golang/prometheus/promhttp
		go get github.com/alecthomas/kingpin/v2

build:
		go build -o bin fohhnnet_exporter
		go build -o bin fohhnnet_exporter/fohhn-cli


run:
		go run main.go

install:
		install -o prometheus -g prometheus bin/fohhnnet_exporter /usr/local/bin
		install -o prometheus -g prometheus bin/fohhn-cli /usr/local/bin/fohhn-cli

