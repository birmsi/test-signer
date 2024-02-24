package service

import (
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/birmsi/test-signer/internal/signatures/api/requests"
	"github.com/birmsi/test-signer/internal/signatures/repository"
	"github.com/birmsi/test-signer/internal/signatures/repository/models"
	"github.com/birmsi/test-signer/internal/signatures/service/domain"
)

type MockSignaturesRepository struct {
	UserSignature domain.UserSignature
	Err           error
}

func (m *MockSignaturesRepository) SignAnswers(signature models.UserSignature) error {
	return m.Err
}

func (m *MockSignaturesRepository) GetUserSignature(userID string, signature []byte) (domain.UserSignature, error) {
	return m.UserSignature, m.Err
}

func MockSignerService(logger slog.Logger, repository repository.SignaturesRepositoryInterface) SignaturesService {
	return SignaturesService{
		logger:               logger,
		signaturesRepository: repository,
	}
}

func TestSign(t *testing.T) {
	// Mock dependencies
	slogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	mockRepository := &MockSignaturesRepository{}
	mockService := MockSignerService(*slogger, mockRepository)

	signRequest := requests.PostSignature{
		Jwt:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOiJqb3NlIn0.ZkhIJqEf0Uy_tIEcap9TtonWLXkNs96LW_7z21fKqhM",
		Answers: []string{"answer1", "answer2"},
	}

	// Test case: successful sign
	signature, err := mockService.Sign(signRequest)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if signature == nil {
		t.Error("Expected non-nil signature, got nil")
	}

	// Test case: error in getting userID from JWT
	signRequest.Jwt = "invalid_jwt"
	_, err = mockService.Sign(signRequest)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Test case: error in signing answers
	mockRepository.Err = errors.New("repository error")
	signRequest.Jwt = "mock_jwt"
	_, err = mockService.Sign(signRequest)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestVerify(t *testing.T) {
	// Mock dependencies
	slogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	mockRepository := &MockSignaturesRepository{}
	mockService := MockSignerService(*slogger, mockRepository)

	// Mock data
	verifyRequest := requests.PostVerifySignature{
		Jwt:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOiJqb3NlIn0.ZkhIJqEf0Uy_tIEcap9TtonWLXkNs96LW_7z21fKqhM",
		Signature: []byte("mock_signature"),
	}

	userSignature, err := mockService.Verify(verifyRequest)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if userSignature.UserID != "" {
		t.Errorf("Expected empty UserID, got %s", userSignature.UserID)
	}

	// Test case: error in getting userID from JWT
	verifyRequest.Jwt = "invalid_jwt"
	_, err = mockService.Verify(verifyRequest)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Test case: error in getting user signature from repository
	mockRepository.Err = errors.New("repository error")
	verifyRequest.Jwt = "mock_jwt"
	_, err = mockService.Verify(verifyRequest)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
