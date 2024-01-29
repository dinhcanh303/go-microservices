package sharedkernel

import "fmt"

type StatusPost int32

const (
	StatusPublic StatusPost = iota + 1
	StatusPrivate
	StatusCompany
	StatusFriends
)

func (s StatusPost) String() string {
	return fmt.Sprintf("%d", int32(s))
}

// func (e StatusPost) String() string {
// 	return []string{
// 		"Public",
// 		"Private",
// 		"Company",
// 		"Friends",
// 	}[e]
// }

type RoleGroupMember int32

const (
	OWNER RoleGroupMember = iota + 1
	ADMIN
	MODERATOR
	USER
)

func (r RoleGroupMember) String() string {
	return fmt.Sprintf("%d", int32(r))
}

type StatusGroupMember int32

const (
	StatusPending StatusGroupMember = iota + 1
	StatusAccepted
	StatusBlocked
)

func (e StatusGroupMember) String() string {
	return fmt.Sprintf("%d", int32(e))
}
