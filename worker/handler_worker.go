package worker

import (
	"context"
	"distributedcharge/accessor"
	"distributedcharge/element"
	"fmt"
)

type HandleWorker struct {
	eventno string
	client  *accessor.ChargeRpc
	trans   *element.Transaction
	dao     accessor.Dao
	done    chan bool
	failed  chan bool
}

func NewHandleWorker(
	dao accessor.Dao,
	trans *element.Transaction,
	client *accessor.ChargeRpc,
	eventno string) *HandleWorker {

	w := &HandleWorker{
		client: client,
		trans:  trans,
		dao:    dao,
		done:   make(chan bool, 1),
		failed: make(chan bool, 1),
	}
	return w
}

func (w *HandleWorker) Start(ctx context.Context) {
	err := w.Execute(w.trans.Nodetask)
	if err != nil {
		w.failed <- true
		w.client.ReportResult(w.eventno, false)
	} else {
		w.client.ReportResult(w.eventno, true)
	}

	select {
	case <-w.done:
		w.Commit()
	case <-w.failed:
		w.RollBack()
	case <-ctx.Done():
		w.RollBack()
	}

}

func (w *HandleWorker) Execute(stmts []element.Stmts) error {
	/*-----------------------
		lock line
		select old value and eventno
		update new value
	------------------------*/
	for _, s := range stmts {
		// lock 失败 重试超时 放弃
		w.dao.LockAccount(s.Account)
		balance := w.dao.GetBalance(s.Account)
		var newbalance int32
		switch s.OP {
		case element.Add:
			newbalance = balance + s.Num
		case element.Minux:
			if (balance - s.Num) < 0 {
				return fmt.Errorf("%s:%d Insufficient balance", w.eventno, s.Account)
			}
			newbalance = balance - s.Num
		}
		w.dao.SetTempBalance(s.Account, newbalance, w.eventno)
	}
	return nil
}

func (w *HandleWorker) RollBack() {
	/*-----------------------
		clear temp
		unlock line
	------------------------*/
	for _, s := range w.trans.Nodetask {
		w.dao.ClearTempBalance(s.Account)
		w.dao.UnlockAccout(s.Account)
	}
}

func (w *HandleWorker) Commit() {
	/*-----------------------
		swap balance
		clear temp
		unlock line
	------------------------*/
	for _, s := range w.trans.Nodetask {
		w.dao.SwapBalance(s.Account)
		w.dao.ClearTempBalance(s.Account)
		w.dao.UnlockAccout(s.Account)
	}
}

func (w *HandleWorker) TaskDone() {
	w.done <- true
}

func (w *HandleWorker) TaskFailed() {
	w.failed <- true
}

func (w *HandleWorker) GetType() WorkerType {
	return HandleType
}

func (w *HandleWorker) GetEventno() string {
	return w.eventno
}
