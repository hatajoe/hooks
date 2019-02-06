package github

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// VerifyMiddleware verify payload
// If webhook secret token is set, you can use this as middleware of handlers
type VerifyMiddleware struct {
	secret string
}

// NewVerifyMiddleware returns the VerifyMiddleware instance
// The first argument `secret` is GitHub webhook secret token
func NewVerifyMiddleware(secret string) *VerifyMiddleware {
	return &VerifyMiddleware{
		secret: secret,
	}
}

// Verify verify the payload signature and call next handler
func (m VerifyMiddleware) Verify(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Read request body failed: %v", err), http.StatusBadRequest)
			return
		}

		if err := verifySignature([]byte(m.secret), r.Header.Get("X-Hub-Signature"), body); err != nil {
			http.Error(w, fmt.Sprintf("Verify signature failed: %v", err), http.StatusUnauthorized)
			return
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		handler.ServeHTTP(w, r)
	})
}

// ref: https://gist.github.com/rjz/b51dc03061dbcff1c521
func verifySignature(secret []byte, signature string, body []byte) error {

	const signaturePrefix = "sha1="
	const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

	if len(signature) != signatureLength {
		return fmt.Errorf("invalid signature length: sig=%s, len=%d", signature, len(signature))
	}

	if !strings.HasPrefix(signature, signaturePrefix) {
		return fmt.Errorf("signature doesn't have prefix")
	}

	actual := make([]byte, 20)
	hex.Decode(actual, []byte(signature[5:]))

	if !hmac.Equal(signBody(secret, body), actual) {
		return fmt.Errorf("signed body is not equal to signature")
	}

	return nil
}

func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	return []byte(computed.Sum(nil))
}
