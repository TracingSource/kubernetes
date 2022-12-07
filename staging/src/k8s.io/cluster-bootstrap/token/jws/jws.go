package bootstrap

import (
	"fmt"
	"strings"

	jose "gopkg.in/square/go-jose.v2"
)

// ComputeDetachedSignature takes content and token details and computes a detached
// JWS signature.  This is described in Appendix F of RFC 7515.  Basically, this
// is a regular JWS with the content part of the signature elided.
func ComputeDetachedSignature(content, tokenID, tokenSecret string) (string, error) {
	jwk := &jose.JSONWebKey{
		Key:   []byte(tokenSecret),
		KeyID: tokenID,
	}

	opts := &jose.SignerOptions{
		// Since this is a symmetric key, go-jose doesn't automatically include
		// the KeyID as part of the protected header. We have to pass it here
		// explicitly.
		ExtraHeaders: map[jose.HeaderKey]interface{}{
			"kid": tokenID,
		},
	}

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: jwk}, opts)
	if err != nil {
		return "", fmt.Errorf("can't make a HS256 signer from the given token: %v", err)
	}

	jws, err := signer.Sign([]byte(content))
	if err != nil {
		return "", fmt.Errorf("can't HS256-sign the given token: %v", err)
	}

	fullSig, err := jws.CompactSerialize()
	if err != nil {
		return "", fmt.Errorf("can't serialize the given token: %v", err)
	}
	return stripContent(fullSig)
}

// stripContent will remove the content part of a compact JWS
//
// The `go-jose` library doesn't support generating signatures with "detached"
// content. To make up for this we take the full compact signature, break it
// apart and put it back together without the content section.
func stripContent(fullSig string) (string, error) {
	parts := strings.Split(fullSig, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("compact JWS format must have three parts")
	}

	return parts[0] + ".." + parts[2], nil
}

// DetachedTokenIsValid checks whether a given detached JWS-encoded token matches JWS output of the given content and token
func DetachedTokenIsValid(detachedToken, content, tokenID, tokenSecret string) bool {
	newToken, err := ComputeDetachedSignature(content, tokenID, tokenSecret)
	if err != nil {
		return false
	}
	return detachedToken == newToken
}
