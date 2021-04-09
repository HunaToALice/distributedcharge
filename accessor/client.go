package accessor

import (
	"distributedcharge/element"
)

type ChargeRpc struct {
}

func (s *ChargeRpc) Requst(*element.ChargeRequest) {

}

func (s *ChargeRpc) ChargePart(uuid string, t *element.Transaction) {
}

func (s *ChargeRpc) Commit(eventno string, uuid string, iscommit bool) {
}

func (s *ChargeRpc) ReportResult(eventno string, result bool) {
}
