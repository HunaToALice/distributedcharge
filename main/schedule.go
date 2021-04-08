package main

import (
	"context"
	"distributedcharge/worker"
	"sync"
	"time"
)

type Schedule struct {
	workerPool chan worker.Worker
	cb         *ControlWorkerBuffer
	hb         *HandleWorkerBuffer
}

type ControlWorkerBuffer struct {
	sync.Mutex
	controlWorkMap map[string]*worker.ControlWorker
}

type HandleWorkerBuffer struct {
	sync.Mutex
	handleWorkMap map[string]*worker.HandleWorker
}

func (s *Schedule) Start() {
	for i := 0; i < 100; i++ {
		go s.DoTask()
	}
}

func (s *Schedule) DoTask() {
	for w := range s.workerPool {
		ctx, _ := context.WithTimeout(context.TODO(), time.Second*3)
		w.Start(ctx)
		s.Eliminate(w.GetEventno())
	}
}

func (s *Schedule) Eliminate(eventno string) {

}

func (b *ControlWorkerBuffer) ControlReport(eventno string, result bool) {
	b.Lock()
	defer b.Unlock()
	w := b.controlWorkMap[eventno]
	if result {
		w.TaskDone()
	} else {
		w.TaskFailed()
	}
}

func (b *HandleWorkerBuffer) HandlerReport(eventno string, commit bool) {
	b.Lock()
	defer b.Unlock()
	w := b.handleWorkMap[eventno]
	if commit {
		w.Commit()
	} else {
		w.RollBack()
	}

}
