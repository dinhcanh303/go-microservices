package token

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

type JWTMaker struct{}

// CreateToken implements Maker.
func (maker *JWTMaker) CreateToken(payload *Payload, secretKey string) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(secretKey))
	return token, err
}

// VerifyToken implements Maker.
func (maker *JWTMaker) VerifyToken(token, secretKey string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(validationErr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}

var _ JWT = (*JWTMaker)(nil)

func NewJWTMaker() JWT {
	return &JWTMaker{}
}
