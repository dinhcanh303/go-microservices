package sharedkernel

import "fmt"

type Status int32

const (
	StatusPrivate Status = iota
	StatusPublic
)

func (e Status) String() string {
	return fmt.Sprintf("%d", int(e))
}
