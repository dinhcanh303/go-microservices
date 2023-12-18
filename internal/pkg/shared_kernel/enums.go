package sharedkernel

type Status int32

const (
	StatusPublic Status = iota
	StatusPrivate
	StatusCompany
	StatusFriends
)

func (e Status) Int32() int32 {
	return int32(e)
}
func (e Status) String() string {
	return []string{
		"Public",
		"Private",
		"Company",
		"Friends",
	}[e]
}
