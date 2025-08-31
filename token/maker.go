package token

import "time"

type TokenMaker interface {
	CreateToken(email string, duration time.Duration) (string, *Payload, error)

	ValidateToken(token string) (*Payload, error)
}
