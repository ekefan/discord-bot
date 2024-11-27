package middleware

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/ekefan/discord-bot/util"
	"github.com/stretchr/testify/require"
)

type discordRequestValues struct {
	signature string
	timestamp string
	message   []byte
}

var config util.EnvConfig
var reqValues discordRequestValues

func TestMain(m *testing.M) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		slog.Error("can't run test for middleware", "details", err.Error())
		os.Exit(1)
	}

	timestamp := time.Now().Unix()
	body := "a random message for testing"
	msg := []byte(fmt.Sprintf("%v%v", timestamp, body))
	signature := ed25519.Sign(privateKey, msg)
	signatureEncoded := hex.EncodeToString(signature)
	publicKeyEncoded := hex.EncodeToString(publicKey)

	config = util.EnvConfig{
		PublicKey: publicKeyEncoded,
	}
	reqValues = discordRequestValues{
		signature: signatureEncoded,
		timestamp: fmt.Sprintf("%v", timestamp),
		message:   []byte(body),
	}

	os.Exit(m.Run())
}

func TestVerifySignature(t *testing.T) {

	testCases := []struct {
		name         string
		signature    string
		timestamp    string
		pubKeyConfig *util.EnvConfig
		reqBody      []byte
		expectedErr  error
	}{
		{
			name:         "ideal case",
			signature:    reqValues.signature,
			timestamp:    reqValues.timestamp,
			pubKeyConfig: &config,
			reqBody:      reqValues.message,
			expectedErr:  nil,
		}, {
			name:         "missing header case 1",
			signature:    reqValues.signature,
			timestamp:    "",
			pubKeyConfig: &config,
			reqBody:      reqValues.message,
			expectedErr:  ErrorMissingHeaderValues,
		}, {
			name:         "missing header case 2",
			signature:    "",
			timestamp:    reqValues.timestamp,
			pubKeyConfig: &config,
			reqBody:      reqValues.message,
			expectedErr:  ErrorMissingHeaderValues,
		}, {
			name:         "invalid header signature",
			signature:    "temperedHeaderWithNonHexEncoding",
			timestamp:    reqValues.timestamp,
			pubKeyConfig: &config,
			reqBody:      reqValues.message,
			expectedErr:  ErrDecodingSignature,
		}, {
			name:      "invalid public key",
			signature: reqValues.signature,
			timestamp: reqValues.timestamp,
			pubKeyConfig: &util.EnvConfig{
				PublicKey: "b1511a9905884e07ff7acdc6a6d6129c87fd52c32f41e77b224274f7c0df",
			},
			reqBody:     reqValues.message,
			expectedErr: ErrInvalidPublicKey,
		}, {
			name:      "invalid public key: non hex encoded key",
			signature: reqValues.signature,
			timestamp: reqValues.timestamp,
			pubKeyConfig: &util.EnvConfig{
				PublicKey: "nonHexEncodedKey",
			},
			reqBody:     reqValues.message,
			expectedErr: ErrDecodingPubKey,
		}, {
			name:         "tempered body",
			signature:    reqValues.signature,
			timestamp:    reqValues.timestamp,
			pubKeyConfig: &config,
			reqBody:      []byte("tempered body"),
			expectedErr:  ErrVerifySignature,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// create an io.Reader from the req body
			body := bytes.NewReader(tc.reqBody)
			// mock a http request to the interactions endpoint
			r := httptest.NewRequest(http.MethodPost, "/interactions", body)
			r.Header.Set("X-Signature-Ed25519", tc.signature)
			r.Header.Set("X-Signature-Timestamp", tc.timestamp)

			// mock a response writer
			w := httptest.NewRecorder()

			// call the verifySignature function
			err := verifySignature(w, r, tc.pubKeyConfig)
			require.Equal(t, tc.expectedErr, err)
		})
	}
}
