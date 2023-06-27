package api

type GConnfd interface {
	ReadMsg()
	HandleMsg()
	PackMsg()
	SendMsg()
	Closefd()
}
