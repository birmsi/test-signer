package domain

import "time"

type SignatureHash struct {
	Answers   []string  `json:"answers"`
	UserID    string    `json:"userID"`
	Timestamp time.Time `json:"timestamp"`
}

func NewSignatureHash(userID string, answers []string, timestamp time.Time) SignatureHash {
	return SignatureHash{
		Answers:   answers,
		UserID:    userID,
		Timestamp: timestamp,
	}
}
