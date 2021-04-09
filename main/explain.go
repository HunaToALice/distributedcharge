package main

import (
	"distributedcharge/accessor"
	"distributedcharge/element"
	"distributedcharge/worker"
)

type Explainer struct {
	s    *Schedule
	crpc *accessor.ChargeRpc
	dao  accessor.Dao
}

/*
	trans outer transaction to inner transaction and conter handler
*/

func (e *Explainer) ExplainChargeRequest(c *element.ChargeRequest, eventno string) *worker.ControlWorker {
	cw := worker.NewControlWorker(eventno, e.crpc)
	// explain
	return cw
}

func (e *Explainer) ExplainTransaction(t *element.Transaction, eventno string) *worker.HandleWorker {
	hw := worker.NewHandleWorker(e.dao, t, e.crpc, eventno)
	return hw
}
