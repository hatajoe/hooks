package github

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifyCorrect(t *testing.T) {

	secret := "test"
	payload := []byte("{}")

	signature := makeHMAC(payload, secret)

	req := httptest.NewRequest("POST", "/webhooks", nil)
	req.Header.Set("X-Hub-Signature", signature)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	verifier := &VerifyMiddleware{
		secret: secret,
	}

	handler := verifier.Verify(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
	}))

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Error("verify sigunature failed")
	}
}

func TestVerifyFailure(t *testing.T) {

	secret := "test"
	payload := []byte("{}")

	signature := makeHMAC(payload, secret)

	req := httptest.NewRequest("POST", "/webhooks", nil)
	req.Header.Set("X-Hub-Signature", signature)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	verifier := &VerifyMiddleware{
		secret: "dummy secret",
	}

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
