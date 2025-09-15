package token

import "time"

// maker is an interface for managing tokens
type Maker interface {
	// createToken creates a new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)

	// check if token is valid
	VerifyToken(token string) (*Payload, error )
}