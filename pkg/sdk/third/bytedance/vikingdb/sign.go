package vikingdb

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"
)

func sign(accessKey, accessSecret, service, region, host, method, path string, params map[string]string, body []byte) (headers map[string]string) {
	headers = make(map[string]string, 4)

	// Build canonical query string (Volcengine spec requires space as %20)
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	canonicalQuery := strings.ReplaceAll(query.Encode(), "+", "%20")

	// Prepare basic headers
	now := time.Now().UTC()
	timestamp := now.Format("20060102T150405Z")
	date := timestamp[:8]

	contentType := "application/json"
	payloadHash := sha256.Sum256(body)
	payloadHex := hex.EncodeToString(payloadHash[:])

	headers["Content-Type"] = contentType
	headers["X-Date"] = timestamp
	headers["X-Content-Sha256"] = payloadHex

	// Canonical headers/signed headers
	signedHeaders := []string{"host", "x-date", "x-content-sha256", "content-type"}
	canonicalHeaderLines := []string{
		"host:" + host,
		"x-date:" + strings.TrimSpace(timestamp),
		"x-content-sha256:" + payloadHex,
		"content-type:" + strings.TrimSpace(contentType),
	}
	canonicalHeaders := strings.Join(canonicalHeaderLines, "\n") + "\n"

	// Canonical request
	canonicalRequest := strings.Join([]string{
		strings.ToUpper(method),
		path,
		canonicalQuery,
		canonicalHeaders,
		strings.Join(signedHeaders, ";"),
		payloadHex,
	}, "\n")

	hashCanonical := sha256.Sum256([]byte(canonicalRequest))
	hashedCanonicalHex := hex.EncodeToString(hashCanonical[:])

	// Signing key derivation
	hmacSHA256 := func(key []byte, data string) []byte {
		mac := hmac.New(sha256.New, key)
		mac.Write([]byte(data))
		return mac.Sum(nil)
	}

	credentialScope := fmt.Sprintf("%s/%s/%s/request", date, region, service)
	signingKey := hmacSHA256([]byte(accessSecret), date)
	signingKey = hmacSHA256(signingKey, region)
	signingKey = hmacSHA256(signingKey, service)
	signingKey = hmacSHA256(signingKey, "request")

	// String to sign & signature
	stringToSign := strings.Join([]string{
		"HMAC-SHA256",
		timestamp,
		credentialScope,
		hashedCanonicalHex,
	}, "\n")
	signature := hex.EncodeToString(hmacSHA256(signingKey, stringToSign))

	// Authorization header
	authorization := fmt.Sprintf(
		"HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		accessKey,
		credentialScope,
		strings.Join(signedHeaders, ";"),
		signature,
	)
	headers["Authorization"] = authorization

	return headers
}
