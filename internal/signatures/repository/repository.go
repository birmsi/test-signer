package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/birmsi/test-signer/internal/signatures/repository/models"
	"github.com/birmsi/test-signer/internal/signatures/service/domain"
	"github.com/lib/pq"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type SignaturesRepositoryInterface interface {
	SignAnswers(userSignature models.UserSignature) error
	GetUserSignature(userID string, signature []byte) (domain.UserSignature, error)
}

type SignaturesRepository struct {
	logger slog.Logger
	db     *sql.DB
}

func NewSignaturesRepository(logger slog.Logger, db *sql.DB) SignaturesRepository {
	return SignaturesRepository{
		logger: logger,
		db:     db,
	}
}

func (sr SignaturesRepository) SignAnswers(userSignature models.UserSignature) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt, err := sr.db.PrepareContext(ctx, "INSERT INTO user_signatures(user_id, signature, answers,hash_timestamp) VALUES($1, $2, $3, $4)")
	if err != nil {
		return nil
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		userSignature.UserID,
		userSignature.Signature,
		pq.Array(userSignature.Answers),
		userSignature.HashTimestamp)

	return err
}

func (sr SignaturesRepository) GetUserSignature(userID string, signature []byte) (domain.UserSignature, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt, err := sr.db.PrepareContext(ctx, "SELECT user_id, signature, answers,hash_timestamp FROM user_signatures WHERE user_id = $1 AND signature = $2")
	if err != nil {
		return domain.UserSignature{}, err
	}
	defer stmt.Close()

	var userSignature models.UserSignature

	err = stmt.QueryRowContext(ctx, userID, signature).Scan(
		&userSignature.UserID,
		&userSignature.Signature,
		pq.Array(&userSignature.Answers),
		&userSignature.HashTimestamp)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return domain.UserSignature{}, ErrRecordNotFound
		default:
			return domain.UserSignature{}, err
		}
	}

	return userSignature.ToDomain(), nil

}
