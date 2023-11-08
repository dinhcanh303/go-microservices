package token

type JWT interface {
	CreateToken(payload *Payload, secretKey string) (string, error)
	VerifyToken(token, secretKey string) (*Payload, error)
}
