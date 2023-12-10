// Package crypto provides a function to encrypt and decrypt data using AES block cipher.
package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func EncryptAES(key []byte, buffer []byte) ([]byte, error) {
	// pad buffer to block size
	if len(buffer)%aes.BlockSize != 0 {
		paddingLen := aes.BlockSize - (len(buffer) % aes.BlockSize)
		padding := bytes.Repeat([]byte{0x00}, paddingLen)
		buffer = append(buffer, padding...)
	}

	// create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure.
	// Therefore, it's common to include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(buffer))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], buffer)

	// return encoded bytes
	return ciphertext, nil
}

func DecryptAES(key []byte, ciphertextBytes []byte) ([]byte, error) {
	// create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// allocate space for deciphered data
	buffer := make([]byte, len(ciphertextBytes))

	iv := ciphertextBytes[:aes.BlockSize]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(buffer, ciphertextBytes)

	return bytes.TrimRight(buffer[aes.BlockSize:], "\x00"), nil
}
