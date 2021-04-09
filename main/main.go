package main

import (
	"distributedcharge/accessor"
	"flag"
)

var (
	uuid        = flag.String("uuid", "", "")
	lasteventno = flag.Uint64("lasteventno", 0, "")
)

func main() {
	ep := NewEventProducer(*uuid, *lasteventno)
	schedule := NewSchedule()
	explainer := &Explainer{
		s:    schedule,
		crpc: &accessor.ChargeRpc{},
		// init dao
	}
	cps := &ChargeRpcServer{
		ep:        ep,
		explainer: explainer,
		schedule:  schedule,
	}
	cps.Init()
}
