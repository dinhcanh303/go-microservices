package sharedkernel

type GetNotiOptions struct {
	Unread bool
	Read   bool
	Limit  int
	Offset int
}
