package consts

type meta int

const (
	// Attempts попытки
	Attempts = meta(0)
	// Retry повторить попытку
	Retry = meta(1)
)
