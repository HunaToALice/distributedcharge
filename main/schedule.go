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

func NewSchedule() *Schedule {
	s := &Schedule{
		workerPool: make(chan worker.Worker, 10000),
		cb:         &ControlWorkerBuffer{controlWorkMap: make(map[string]*worker.ControlWorker)},
		hb:         &HandleWorkerBuffer{handleWorkMap: make(map[string]*worker.HandleWorker)},
	}
	return s
}

func (s *Schedule) Start() {
	for i := 0; i < 100; i++ {
		go s.DoTask()
	}
}

func (s *Schedule) PutControlWorker(eventno string, w *worker.ControlWorker) {
	s.cb.AddWorker(eventno, w)
	s.workerPool <- w
}

func (s *Schedule) PutHandleWorker(eventno string, w *worker.HandleWorker) {
	s.hb.AddWorker(eventno, w)
	s.workerPool <- w
}

func (s *Schedule) DoTask() {
	for w := range s.workerPool {
		ctx, _ := context.WithTimeout(context.TODO(), time.Second*3)
		w.Start(ctx)
		s.Eliminate(w.GetEventno(), w)
	}
}

func (s *Schedule) Eliminate(eventno string, w worker.Worker) {
	switch w.GetType() {
	case worker.ControlType:
		s.cb.DeleteWorker(eventno)
	case worker.HandleType:
		s.hb.DeleteWorker(eventno)
	}
}

func (s *Schedule) NotifyCw(eventno string, result bool) {
	s.cb.ControlReport(eventno, result)
}

func (s *Schedule) NotifyHw(eventno string, iscommit bool) {
	s.hb.HandlerReport(eventno, iscommit)
}

func (b *ControlWorkerBuffer) AddWorker(eventno string, w *worker.ControlWorker) {
	b.Lock()
	defer b.Unlock()
	b.controlWorkMap[eventno] = w
}

func (b *ControlWorkerBuffer) DeleteWorker(eventno string) {
	b.Lock()
	defer b.Unlock()
	delete(b.controlWorkMap, eventno)
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

func (b *HandleWorkerBuffer) AddWorker(eventno string, w *worker.HandleWorker) {
	b.Lock()
	defer b.Unlock()
	b.handleWorkMap[eventno] = w
}

func (b *HandleWorkerBuffer) DeleteWorker(eventno string) {
	b.Lock()
	defer b.Unlock()
	delete(b.handleWorkMap, eventno)
}
