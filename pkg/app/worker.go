package app

type Worker interface {
	Start() error
	Stop() error
}
