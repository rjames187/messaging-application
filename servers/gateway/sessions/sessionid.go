package sessions

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

func createSessionToken(secret string, IDLength int) (string, string, error) {
	id, err := generateRandomBytes(IDLength); if err != nil {
		return "", "", err
	}
	secretBytes, err := base64.URLEncoding.DecodeString(secret)
	if err != nil {
		return "", "", err
	}
	signature, err := signID(id, secretBytes); if err != nil {
		return "", "", err
	}
	token := append(id, signature...)
	return encode(token), encode(id), nil
}

func validToken(sessionToken string, secret string, IDLength int) (bool, string, error) {
	tokenBytes, err := base64.URLEncoding.DecodeString(sessionToken)
	if err != nil {
		return false, "", err
	}
	sessionID := tokenBytes[:IDLength]
	signature := tokenBytes[IDLength:]
	secretBytes, err := base64.URLEncoding.DecodeString(secret)
	if err != nil {
		return false, "", err
	}
	signatureToCompare, err := signID(sessionID, secretBytes)
	if err != nil {
		return false, "", err
	}
	if !hmac.Equal(signature, signatureToCompare) {
		return false, "", nil
	}
	return true, encode(sessionID), nil
}

func generateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes); if err != nil {
		return nil, err
	}
	return bytes, nil
}

func signID(id []byte, secret []byte) ([]byte, error) {
	h := hmac.New(sha256.New, secret)
	_, err := h.Write(id); if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func encode(token []byte) string {
	return base64.URLEncoding.EncodeToString(token)
}