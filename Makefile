config:
        go get github.com/go-kit/kit/log
		go get github.com/go-kit/kit/log/level
		go get github.com/prometheus/common/version
		go get github.com/prometheus/client_golang/prometheus
		go get github.com/prometheus/client_golang/prometheus/promhttp
		go get github.com/alecthomas/kingpin/v2

build:
		go build fohhnnet_exporter
		go build fohhnnet_exporter/fohhn-cli


run:
		go run main.go

install:
		install -o prometheus -g prometheus fohhnnet_exporter /usr/local/bin
		install -o prometheus -g prometheus fohhn-cli /usr/local/bin
