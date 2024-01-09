package sessions

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

func createSessionID(secret []byte) (string, error) {
	bytes, err := generateRandomBytes(32); if err != nil {
		return "", err
	}
	signature, err := signID(bytes, secret); if err != nil {
		return "", err
	}
	id := encodeSignature(signature)
	return id, nil
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

func encodeSignature(signature []byte) string {
	return base64.URLEncoding.EncodeToString(signature)
}