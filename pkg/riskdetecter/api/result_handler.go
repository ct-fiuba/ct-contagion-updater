package api

type ResultHandler interface {
	Handle(r Result)
}
