package main

import (
	"fohhnnet_exporter/FohhnNet"
	"github.com/alecthomas/kingpin/v2"
)

var (
	debugOn = kingpin.Flag("debug", "Get debug informations").Bool()
	id      = kingpin.Flag("id", "Get debug informations").Int8()
	scan    = kingpin.Flag("scan", "Scan for devices").Bool()
	all     = kingpin.Flag("all", "Show data from all devices").Bool()
	host    = kingpin.Arg("host", "Host or IP of target device").Required().String()
	port    = kingpin.Flag("port", "Port of target device").Default("4001").Int()
)

func main() {

	kingpin.Parse()

	scanNet()

}

func scanNet() {
	fohhnNetSession, err := FohhnNet.NewFohhnNetTcpSession(*host, *port)
	hasAnswered := false
	_ = hasAnswered

	if err == nil {

		defer FohhnNet.Close(fohhnNetSession) // finally, close the connection

		if fohhnNetSession.IsConnected {

			if *scan {
				res := FohhnNet.ScanFohhnNet(fohhnNetSession, 1, 2)

				FohhnNet.RenderFohhnNetScan(res)

				if *all {
					for _, id := range res {

						walkResult, err := FohhnNet.ScrapeFohhnDevice(fohhnNetSession, id)
						if err == nil {
							FohhnNet.RenderFohhnDevice(walkResult)
						}

					}
				}
			}
			if *id > 0 {

				walkResult, err := FohhnNet.ScrapeFohhnDevice(fohhnNetSession, *id)
				if err == nil {
					FohhnNet.RenderFohhnDevice(walkResult)
				}

			}

		}
	}
}
