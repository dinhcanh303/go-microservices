package sharedkernel

type StatusPost int32

const (
	StatusPublic StatusPost = iota + 1
	StatusPrivate
	StatusCompany
	StatusFriends
)

func (e StatusPost) String() string {
	return []string{
		"public",
		"private",
		"company",
		"friends",
	}[e]
}

type RoleGroupMember int32

const (
	NO_USER RoleGroupMember = iota
	OWNER
	ADMIN
	MODERATOR
	USER
)

func (r RoleGroupMember) String() string {
	return []string{
		"",
		"owner",
		"admin",
		"moderator",
		"user",
	}[r]
}

type StatusGroupMember int32

const (
	StatusPending StatusGroupMember = iota + 1
	StatusAccepted
	StatusBlocked
)

func (r StatusGroupMember) String() string {
	return []string{
		"pending",
		"accepted",
		"blocked",
	}[r]
}
