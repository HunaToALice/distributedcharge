package accessor

type Acceptor struct {
}

func (a *Acceptor) GetRequst() *ChargeRequest {
	return nil
}

func (a *Acceptor) TransDone(eventno string) string {
	return ""
}

func (a *Acceptor) GetTrans() *Transaction {
	return nil
}
