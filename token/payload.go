package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID
	Email     string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("this token has expired")
)

func NewPayload(email string, duration time.Duration) *Payload {
	issuedAt := time.Now()

	return &Payload{
		ID:        uuid.New(),
		Email:     email,
		IssuedAt:  issuedAt,
		ExpiresAt: issuedAt.Add(duration),
	}
}
