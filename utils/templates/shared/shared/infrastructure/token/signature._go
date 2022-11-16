package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

type ClientCredential struct {
	Secret []byte
}

func (r ClientCredential) CreateSignature(payloadInBytes []byte) (string, error) {

	h := hmac.New(sha256.New, r.Secret)

	_, err := h.Write(payloadInBytes)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil

}

func (r ClientCredential) ValidateSignature(signatureGiven string, payloadInBytes []byte) error {

	signatureFromPayload, err := r.CreateSignature(payloadInBytes)
	if err != nil {
		return err
	}

	if signatureGiven != signatureFromPayload {
		return fmt.Errorf("invalid signature")
	}
	return nil
}
