package token

import "time"

// Maker is a token maker.
type Maker interface {
	// CreateToken Make creates a new token.
	CreateToken(email string, duration time.Duration) (string, error)
	// VerifyToken Verify verifies a token.
	VerifyToken(token string) (*Payload, error)
}
