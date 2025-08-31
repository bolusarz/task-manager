package token

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	maker        paseto.V4SymmetricKey
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (TokenMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := PasetoMaker{
		maker:        paseto.NewV4SymmetricKey(),
		symmetricKey: []byte(symmetricKey),
	}

	return &maker, nil
}

func (m *PasetoMaker) CreateToken(email string, duration time.Duration) (string, *Payload, error) {
	payload := NewPayload(email, duration)

	token := paseto.NewToken()

	token.SetExpiration(payload.ExpiresAt)
	token.SetIssuedAt(payload.IssuedAt)
	token.SetSubject(payload.Email)
	token.SetNotBefore(time.Now())

	err := token.Set("payload", payload)
	if err != nil {
		return "", nil, err
	}

	return token.V4Encrypt(m.maker, m.symmetricKey), payload, nil
}

func (m *PasetoMaker) ValidateToken(token string) (*Payload, error) {
	parser := paseto.NewParser()
	parsedToken, err := parser.ParseV4Local(m.maker, token, m.symmetricKey)

	if err != nil {
		return nil, err
	}

	payload := &Payload{}
	err = parsedToken.Get("payload", payload)

	if err != nil {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
