package accessor

type Dao interface {
	LockAccount(account uint64)
	UnlockAccout(account uint64) // try unlock
	GetBalance(account uint64) int32
	SetTempBalance(account uint64, balance int32, eventno string)
	SwapBalance(account uint64)
	ClearTempBalance(account uint64)
}
