package main

import (
	"log"
	"net"
	"time"

	"github.com/jwhited/bgpls"
	flag "github.com/spf13/pflag"
)

var (
	casn     = flag.Int64("collector-asn", 64512, "ASN number of the collector")
	nasn     = flag.Int64("neighbor-asn", 64512, "ASN number of the neighbor")
	routerID = flag.String("router-id", "", "Router ID of the collector")
	address  = flag.String("address", "", "Address of the neighbor")
	holdTime = flag.Int64("hold-time", 30, "Hold time in seconds")
)

func main() {
	flag.Parse()
	collectorConfig := &bgpls.CollectorConfig{
		ASN:             uint32(*casn),
		RouterID:        net.ParseIP(*routerID),
		EventBufferSize: 1024,
	}

	collector, err := bgpls.NewCollector(collectorConfig)
	if err != nil {
		log.Fatal(err)
	}

	neighborConfig := &bgpls.NeighborConfig{
		Address:  net.ParseIP(*address),
		ASN:      uint32(*nasn),
		HoldTime: time.Second * time.Duration(*holdTime),
	}

	err = collector.AddNeighbor(neighborConfig)
	if err != nil {
		log.Fatal(err)
	}

	eventsChan, err := collector.Events()
	if err != nil {
		log.Fatal(err)
	}

	for {
		event := <-eventsChan
		// all Event types can be found in event.go (EventNeighbor**)
		switch e := event.(type) {
		case *bgpls.EventNeighborErr:
			log.Printf("neighbor %s, err: %v", e.Neighbor().Address, e.Err)
		case *bgpls.EventNeighborStateTransition:
			log.Printf("neighbor %s, state transition: %v", e.Neighbor().Address, e.State)
		case *bgpls.EventNeighborNotificationReceived:
			log.Printf("neighbor %s notification message code: %v", e.Neighbor().Address, e.Message.Code)
		case *bgpls.EventNeighborUpdateReceived:
			log.Printf("neighbor %s update message", e.Neighbor().Address)
		}
	}
}
