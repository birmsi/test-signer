package service

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/birmsi/test-signer/internal/helpers"
	"github.com/birmsi/test-signer/internal/signatures/api/requests"
	"github.com/birmsi/test-signer/internal/signatures/repository"
	"github.com/birmsi/test-signer/internal/signatures/repository/models"
	"github.com/birmsi/test-signer/internal/signatures/service/domain"
)

type SignaturesServiceInterface interface {
	Sign(signRequest requests.PostSignature) ([]byte, error)
	Verify(verifyRequest requests.PostVerifySignature) (domain.UserSignature, error)
}

type SignaturesService struct {
	logger               slog.Logger
	signaturesRepository repository.SignaturesRepositoryInterface
}

func NewSignaturesService(logger slog.Logger, signaturesRepository repository.SignaturesRepositoryInterface) SignaturesService {
	return SignaturesService{
		logger:               logger,
		signaturesRepository: signaturesRepository,
	}
}

func (ss SignaturesService) Sign(signRequest requests.PostSignature) ([]byte, error) {

	userID, err := helpers.GetUserIDFromJWT(signRequest.Jwt)
	if err != nil {
		ss.logger.Error(err.Error())
		return nil, err
	}

	hashTimestamp := time.Now()

	signature, err := ss.generateHash(userID, signRequest, hashTimestamp)
	if err != nil {
		ss.logger.Error(err.Error())
		return nil, err
	}

	userSignature := domain.UserSignature{
		UserID:        userID,
		Signature:     signature,
		Answers:       signRequest.Answers,
		HashTimestamp: hashTimestamp,
	}

	if err = ss.signaturesRepository.SignAnswers(models.FromDomain(userSignature)); err != nil {
		ss.logger.Error(err.Error())
		return nil, err
	}

	return signature, nil
}

func (ss SignaturesService) Verify(verifyRequest requests.PostVerifySignature) (domain.UserSignature, error) {
	userID, err := helpers.GetUserIDFromJWT(verifyRequest.Jwt)
	if err != nil {
		ss.logger.Error(err.Error())
		return domain.UserSignature{}, err
	}

	return ss.signaturesRepository.GetUserSignature(userID, verifyRequest.Signature)
}

func (ss SignaturesService) generateHash(userID string, signRequest requests.PostSignature, hashTimestamp time.Time) ([]byte, error) {

	serializedData, err := json.Marshal(domain.NewSignatureHash(userID, signRequest.Answers, hashTimestamp))
	if err != nil {
		ss.logger.Error(err.Error())
		return nil, err
	}

	hash := sha256.Sum256(serializedData)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		ss.logger.Error(err.Error())
		return nil, err
	}

	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
}
