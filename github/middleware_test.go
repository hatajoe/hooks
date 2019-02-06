package github

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifyMiddleware_Verify_Correct(t *testing.T) {
	secret := "test secret"
	payload := []byte("{}")
	signature := makeSignature([]byte(secret), payload)

	req := httptest.NewRequest("POST", "/webhooks", nil)
	req.Header.Set("X-Hub-Signature", signature)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	verifier := NewVerifyMiddleware(secret)

	handler := verifier.Verify(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
	}))

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Error(rec.Body)
	}
}

func TestVerifyMiddleware_Verify_Failure(t *testing.T) {
	secret := "test secret"
	payload := []byte("{}")
	signature := makeSignature([]byte(secret), payload)

	req := httptest.NewRequest("POST", "/webhooks", nil)
	req.Header.Set("X-Hub-Signature", signature)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	verifier := NewVerifyMiddleware("dummy secret")

	handler := verifier.Verify(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
	}))

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("returning HTTP status code is not expected: %v", rec.Code)
	}
}

func makeSignature(secret, payload []byte) string {
	return "sha1=" + hex.EncodeToString(signBody(secret, payload))
}
