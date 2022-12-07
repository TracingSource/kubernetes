package crypto

import (
	"github.com/pkg/errors"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

// CreateRandBytes returns a cryptographically secure slice of random bytes with a given size
func CreateRandBytes(size uint32) ([]byte, error) {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}

// EncryptBytes takes a byte slice of raw data and an encryption key and returns an encrypted byte slice of data.
// The key must be an AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256
func EncryptBytes(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce, err := CreateRandBytes(uint32(gcm.NonceSize()))
	if err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, data, nil), nil
}

// DecryptBytes takes a byte slice of encrypted data and an encryption key and returns a decrypted byte slice of data.
// The key must be an AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256
func DecryptBytes(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("size of data is less than the nonce")
	}

	nonce, out := data[:nonceSize], data[nonceSize:]
	out, err = gcm.Open(nil, nonce, out, nil)
	if err != nil {
		return nil, err
	}
	return out, nil
}
