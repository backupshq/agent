package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

func GetAccessTokenExpiry(token string) (time.Time, error) {
	var data = struct {
		Expiry int64 `json:"exp"`
	}{}

	err := decodeToken(&data, token)

	return time.Unix(data.Expiry, 0), err
}

func GetAccessTokenPrincipalId(token string) (string, error) {
	var data = struct {
		PrincipalId string `json:"sub"`
	}{}

	err := decodeToken(&data, token)

	return data.PrincipalId, err
}

func decodeToken(data interface{}, token string) error {
	pieces := strings.Split(token, ".")
	if len(pieces) != 3 {
		return errors.New("Invalid token format")
	}

	jsonBytes, err := base64.RawURLEncoding.DecodeString(pieces[1])
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonBytes, &data)
}
