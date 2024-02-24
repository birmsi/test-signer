package domain

import "time"

type UserSignature struct {
	UserID        string
	Signature     []byte
	Answers       []string
	HashTimestamp time.Time
}
