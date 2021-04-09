package element

type RequestType uint32

const (
	OuterRe RequestType = iota
	InterRe
)

type ChargeRequest struct {
}

type Transaction struct {
	Eventno  string
	Nodetask []Stmts
}

type Operator uint32

const (
	Add Operator = iota
	Minux
)

type Stmts struct {
	Account uint64
	OP      Operator
	Num     int32
}
