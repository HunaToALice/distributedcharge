package worker

import (
	"context"
	"distributedcharge/accessor"
	"distributedcharge/element"
	"fmt"
	"sync"
)

type WorkerType uint16

const (
	ControlType WorkerType = iota
	HandleType
)

type Worker interface {
	Start(ctx context.Context)
	GetEventno() string
	GetType() WorkerType
}

type ControlWorker struct {
	sync.Mutex
	eventno    string
	client     *accessor.ChargeRpc
	innertrans map[string]*element.Transaction
	tonum      uint
	done       chan bool
	failed     chan bool
}

func NewControlWorker(
	eventno string,
	client *accessor.ChargeRpc) *ControlWorker {
	w := &ControlWorker{
		eventno: eventno,
		client:  client,
		done:    make(chan bool, 1),
		failed:  make(chan bool, 1),
	}
	return w
}

func (c *ControlWorker) Start(ctx context.Context) {
	for uuid, t := range c.innertrans {
		go c.client.ChargePart(uuid, t)
	}
	select {
	case <-c.done:
		// 正常结束
		// 通知commit
		c.Commit()
	case <-c.failed:
		// 节点失败，通知回滚
		fmt.Println("node transaction failed")
		c.RollBack()
	case <-ctx.Done():
		// 超时结束 通知回滚
		fmt.Println("timeout")
		c.RollBack()
	}
}

func (c *ControlWorker) TaskDone() {
	c.Lock()
	defer c.Unlock()
	c.tonum--
	if c.tonum == 0 {
		c.done <- true
	}
}

func (c *ControlWorker) TaskFailed() {
	c.failed <- true
}

func (c *ControlWorker) Commit() {
	for uuid, t := range c.innertrans {
		go c.client.Commit(t.Eventno, uuid, true)
	}
}

func (c *ControlWorker) RollBack() {
	for uuid, t := range c.innertrans {
		go c.client.Commit(t.Eventno, uuid, false)
	}
}

func (c *ControlWorker) GetType() WorkerType {
	return ControlType
}

func (c *ControlWorker) GetEventno() string {
	return c.eventno
}
