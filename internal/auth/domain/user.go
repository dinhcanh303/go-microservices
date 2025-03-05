package domain

import (
	"encoding/json"
	"time"

	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/google/uuid"
)

type User struct {
	sharedkernel.AggregateRoot
	ID          uuid.UUID       `json:"id"`
	Email       string          `json:"email"`
	Password    string          `json:"password"`
	FirstName   string          `json:"first_name"`
	LastName    string          `json:"last_name"`
	FullName    string          `json:"full_name"`
	NickName    string          `json:"nick_name"`
	Role        string          `json:"role"`
	AvatarUrl   string          `json:"avatar_url"`
	ProfileUrl  string          `json:"profile_url"`
	Resigned    bool            `json:"resigned"`
	Gender      bool            `json:"gender"`
	Phone       string          `json:"phone"`
	Address     string          `json:"address"`
	Position    string          `json:"position"`
	DateOfBirth time.Time       `json:"date_of_birth"`
	Settings    json.RawMessage `json:"settings"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
type UserFollow struct {
	Id        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	NickName  string    `json:"nick_name"`
	AvatarUrl string    `json:"avatar_url"`
}
type Settings struct {
	Social Social `json:"social"`
}
type Social struct {
	Post   Post   `json:"post"`
	System System `json:"system"`
}
type Post struct {
	StatusDefault int32 `json:"status_default"`
}
type System struct {
	Theme string `json:"theme"`
}
type UserAuth struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewUser(email string, password string, firstName string, lastName string, avatarUrl string) *User {

	return &User{
		ID:        uuid.New(),
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		FullName:  firstName + " " + lastName,
		Password:  password,
		AvatarUrl: avatarUrl,
	}
}
