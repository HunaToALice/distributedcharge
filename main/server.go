package main

import (
	"distributedcharge/element"
)

type ChargeRpcServer struct {
	ep        *EventProducer
	explainer *Explainer
	schedule  *Schedule
}

func (s *ChargeRpcServer) Requst(crq *element.ChargeRequest) {
	eventno := s.ep.GetEventno()
	cw := s.explainer.ExplainChargeRequest(crq, eventno)
	s.schedule.PutControlWorker(eventno, cw)
}

func (s *ChargeRpcServer) ChargePart(eventno string, t *element.Transaction) {
	hw := s.explainer.ExplainTransaction(t, eventno)
	s.schedule.PutHandleWorker(eventno, hw)
}

func (s *ChargeRpcServer) Commit(eventno string, iscommit bool) {
	s.schedule.NotifyHw(eventno, iscommit)
}

func (s *ChargeRpcServer) ReportResult(eventno string, result bool) {
	s.schedule.NotifyCw(eventno, result)
}

func (s *ChargeRpcServer) Init() {}
