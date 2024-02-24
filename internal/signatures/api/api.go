package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/birmsi/test-signer/internal/helpers"
	"github.com/birmsi/test-signer/internal/signatures/api/requests"
	"github.com/birmsi/test-signer/internal/signatures/service"
)

type Signatures struct {
	logger  slog.Logger
	service service.SignaturesService
}

func NewSignaturesAPI(sLogger slog.Logger, service service.SignaturesService) Signatures {
	return Signatures{
		logger:  sLogger,
		service: service,
	}
}

func (s Signatures) Handlers(mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("%s /test/sign", http.MethodPost), s.postSign)
	mux.HandleFunc(fmt.Sprintf("%s /signature/verify", http.MethodPost), s.postVerify)
}

func (s Signatures) postSign(w http.ResponseWriter, r *http.Request) {

	var signatureRequest requests.PostSignature
	if err := helpers.ReadJSON(w, r, &signatureRequest); err != nil {
		s.logger.Error(err.Error())
		helpers.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	signature, err := s.service.Sign(signatureRequest)
	if err != nil {
		s.logger.Error(err.Error())
		helpers.ServerErrorResponse(w, r, err)
		return
	}

	err = helpers.WriteJSON(w, http.StatusCreated, helpers.JsonEnvelope{"signature": signature}, nil)
	if err != nil {
		s.logger.Error(err.Error())
		helpers.ServerErrorResponse(w, r, err)
	}
}

func (s Signatures) postVerify(w http.ResponseWriter, r *http.Request) {
	var signatureRequest requests.PostVerifySignature
	if err := helpers.ReadJSON(w, r, &signatureRequest); err != nil {
		helpers.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	userSignature, err := s.service.Verify(signatureRequest)
	if err != nil {
		s.logger.Error(err.Error())
		helpers.ServerErrorResponse(w, r, err)
		return
	}
	err = helpers.WriteJSON(
		w,
		http.StatusOK,
		helpers.JsonEnvelope{"answers": userSignature.Answers, "timestamp": userSignature.HashTimestamp},
		nil)

	if err != nil {
		s.logger.Error(err.Error())
		helpers.ServerErrorResponse(w, r, err)
	}
}
