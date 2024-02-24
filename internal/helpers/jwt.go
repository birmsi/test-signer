package helpers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

type Claims struct {
	UserID string `json:"userid"`
}

func GetUserIDFromJWT(jwt string) (string, error) {
	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		return "", errors.New("invalid jwt token format")
	}

	claimsData, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	var claims Claims
	err = json.Unmarshal(claimsData, &claims)
	if err != nil {
		return "", err
	}

	return claims.UserID, nil
}
