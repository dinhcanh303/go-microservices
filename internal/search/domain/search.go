package domain

import "github.com/google/uuid"

type Search struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Description string    `json:"description"`
	FullName    string    `json:"full_name"`
	NickName    string    `json:"nick_name"`
	Phone       string    `json:"phone"`
	Position    string    `json:"position"`
	AvatarUrl   string    `json:"avatar_url"`
	ProfileUrl  string    `json:"profile_url"`
}
