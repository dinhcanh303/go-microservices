package sharedkernel

import "fmt"

type Status int8

const (
	StatusPrivate Status = iota
	StatusPublic
)

func (e Status) String() string {
	return fmt.Sprintf("%d", int(e))
}
