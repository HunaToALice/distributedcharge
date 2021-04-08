package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type NumberUnit struct {
	sync.Mutex
	currentnum uint64
}

func NewNumberUnit(lasteventno uint64) *NumberUnit {
	n := &NumberUnit{}
	n.InitNum(lasteventno)
	return n
}

func (n *NumberUnit) InitNum(lasteventno uint64) {
	n.currentnum = lasteventno
}

func (n *NumberUnit) GetNum() uint64 {
	n.Lock()
	defer n.Unlock()
	n.currentnum++
	return n.currentnum
}

type EventProducer struct {
	uuid string
	n    *NumberUnit
}

func NewEventProducer(uuid string, lasteventno uint64) *EventProducer {
	e := &EventProducer{
		uuid: uuid,
		n:    NewNumberUnit(lasteventno)}
	return e
}

func (e *EventProducer) GetEventno() string {
	return fmt.Sprintf("%s:%d", e.uuid, e.n.GetNum())
}

func (e *EventProducer) EventnoParser(eventno string) (string, uint64, error) {
	ss := strings.Split(eventno, ":")
	if len(ss) != 2 {
		return "", 0, fmt.Errorf("Eventformat error")
	}
	num, err := strconv.ParseUint(ss[1], 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("Eventformat error")
	}
	return ss[0], num, nil
}
