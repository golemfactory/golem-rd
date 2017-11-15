package bootstrap

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/ishbir/elliptic"
)

func TestEncryptDecrypt(t *testing.T) {
	rand.Seed(42)
	privKey1, err := elliptic.GeneratePrivateKey(elliptic.Secp256k1)
	if err != nil {
		t.Fatal(err)
	}
	privKey2, err := elliptic.GeneratePrivateKey(elliptic.Secp256k1)
	if err != nil {
		t.Fatal(err)
	}
	data := []byte("asdf")
	encrytped, err := EncryptPython(privKey1, data, &privKey2.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	decrypted, err := DecryptPython(privKey2, encrytped)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, decrypted) {
		t.Fatalf("Expected %v, got %v", data, decrypted)
	}
}
