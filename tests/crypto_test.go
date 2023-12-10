package tests

import (
	"crypto/aes"
	"dev.hackerman.me/artheon/l7-shared-launcher/crypto"
	"testing"
)

func TestAesEncryptDecrypt(t *testing.T) {
	key := []byte("1234567890123456")
	text := "hello world"
	blockSize := aes.BlockSize
	sourceLen := len(text)
	bufferSize := blockSize + sourceLen + (blockSize - sourceLen%blockSize)
	buffer := make([]byte, bufferSize)

	for i, v := range text {
		buffer[i] = byte(v)
	}

	bytes, err := crypto.EncryptAES(key, buffer)
	if err != nil {
		t.Fatal(err)
	}

	if len(bytes) == 0 {
		t.Fatal("bytes is empty")
	}

	decrypted, err := crypto.DecryptAES(key, bytes)
	if err != nil {
		t.Fatal(err)
	}

	if string(decrypted) != text {
		t.Logf("decrypted: %s, original: %s\n", decrypted, text)
		t.Fatal("decrypted text does not match original text")
	}
}
