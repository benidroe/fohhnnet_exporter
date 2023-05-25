package main

import (
	"fmt"
	"fohhnnet_exporter/FohhnNet"
	"github.com/alecthomas/kingpin/v2"
)

var (
	debugOn   = kingpin.Flag("debug", "Get debug informations").Bool()
	id        = kingpin.Flag("id", "Get debug informations").Int8()
	scan      = kingpin.Flag("scan", "Scan for devices").Bool()
	all       = kingpin.Flag("all", "Show data from all devices").Bool()
	host      = kingpin.Arg("host", "Host or IP of target device").Required().String()
	port      = kingpin.Flag("port", "Port of target device").Default("2101").Int()
	connProto = kingpin.Flag("protocol", "Use tcp or udp").Short('p').Default("udp").String()
)

func main() {

	kingpin.Parse()

	// check if selected protocol is valid
	if *connProto != "tcp" && *connProto != "udp" {

		fmt.Printf("%s is not a valid protocol. Choose tcp or udp.\n", *connProto)
		return
	}
	// Set default port for tcp. 4001 is the RS-485 device server port.
	if *connProto == "tcp" && *port == 2101 {
		*port = 4001
	}

	scanNet()

}

func scanNet() {
	fohhnNetSession, err := FohhnNet.NewFohhnNetSession(*host, *port, *connProto)
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
