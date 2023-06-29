package api

type GServer interface {
	ListenAndServe() error
}