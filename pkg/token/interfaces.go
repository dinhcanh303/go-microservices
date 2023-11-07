package token

import "time"

type Maker interface {
	CreateToken(payload interface{}, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
