package middleware

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/ekefan/discord-bot/util"
)

var (
	ErrVerifySignature      = errors.New("signature, could not be verified")
	ErrDecodingSignature     = errors.New("error decoding the hex signature")
	ErrDecodingPubKey        = errors.New("error decoding the hex public key")
	ErrInvalidPublicKey      = errors.New("environment config, public key is incorrect")
	ErrReadingRequestbody    = errors.New("error reading the request body")
	ErrorMissingHeaderValues = errors.New("signature or timestamp header values are missing")
)

func VerifyDiscordSignature(f http.HandlerFunc, config *util.EnvConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := verifySignature(w, r, config); err != nil {
			slog.Error("couldn't verify discord signature", "details", err.Error())
			return
		}
		f(w, r)
	}
}

// verifySignature reads a signature and timestamp from the request
// header and verifies it based on a the body of the request
func verifySignature(w http.ResponseWriter, r *http.Request, config *util.EnvConfig) error {
	// read discords security headers
	signature := r.Header.Get("X-Signature-Ed25519")
	timestamp := r.Header.Get("X-Signature-Timestamp")

	if signature == "" || timestamp == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ErrorMissingHeaderValues
	}

	// read request body to be able to get request message
	body, err := io.ReadAll((r.Body))
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return ErrReadingRequestbody
	}
	// restore request body to the stream
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	// get message
	message := append([]byte(timestamp), body...)

	// Decode the keys and verify the signature

	pubKeyBytes, err := hex.DecodeString(config.PublicKey)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		slog.Error("error decoding public key", "details", err.Error())
		return ErrDecodingPubKey
	}
	if len(pubKeyBytes) != ed25519.PublicKeySize {
		slog.Error("incorrect discord public key size used")
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return ErrInvalidPublicKey
	}
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		slog.Error("error decoding discord signature", "details", err.Error())
		return ErrDecodingSignature
	}

	if !ed25519.Verify(pubKeyBytes, message, sigBytes) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return ErrVerifySignature
	}
	return nil
}
