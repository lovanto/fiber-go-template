package models

import "github.com/google/uuid"

type Renew struct {
	RefreshToken string `json:"refresh_token"`
}

type Tokens struct {
	Access  string
	Refresh string
}

type TokenMetadata struct {
	UserID      uuid.UUID
	Credentials map[string]bool
	Expires     int64
}
