package models

import (
	"time"

	"github.com/birmsi/test-signer/internal/signatures/service/domain"
)

type UserSignature struct {
	UserID        string
	Signature     []byte
	Answers       []string
	HashTimestamp time.Time
}

func (us UserSignature) ToDomain() domain.UserSignature {
	return domain.UserSignature{
		UserID:        us.UserID,
		Signature:     us.Signature,
		Answers:       us.Answers,
		HashTimestamp: us.HashTimestamp,
	}
}

func FromDomain(domainSignature domain.UserSignature) UserSignature {
	return UserSignature{
		UserID:        domainSignature.UserID,
		Signature:     domainSignature.Signature,
		Answers:       domainSignature.Answers,
		HashTimestamp: domainSignature.HashTimestamp,
	}
}
