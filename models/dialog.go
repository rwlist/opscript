package models

type DialogStatus int

const (
	Std DialogStatus = iota
	WaitSrc
)

type Dialog struct {
	Ns     *Namespace
	Status DialogStatus
}
