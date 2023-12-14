package sharedkernel

type Payload struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Role      string `json:"role"`
	AvatarUrl string `json:"avatar_url"`
}
