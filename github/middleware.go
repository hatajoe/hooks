package github

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
)

// VerifyMiddleware verify payload
// If webhook secret token is set, you can use this as middleware of handlers
type VerifyMiddleware struct {
	secret string
}

// Verify verify the payload signature and call next handler
func (m VerifyMiddleware) Verify(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Read request body failed: %v", err), http.StatusBadRequest)
			return
		}

		signature := r.Header.Get("X-Hub-Signature")
		if signature != makeHMAC(body, m.secret) {
			http.Error(w, fmt.Sprintf("Verify signature failed: %v", err), http.StatusUnauthorized)
			return
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		handler.ServeHTTP(w, r)
	})
}

func makeHMAC(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}
