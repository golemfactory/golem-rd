package bootstrap

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math/rand"

	"github.com/ishbir/elliptic"
)

func EncryptPython(key *elliptic.PrivateKey, data []byte, pubkey *elliptic.PublicKey) ([]byte, error) {
	ecdhKey, err := key.GetRawECDHKey(pubkey, 32)
	if err != nil {
		return nil, errors.New("failed to get ECDH key: " + err.Error())
	}

	cipher, err := elliptic.GetCipherByName("aes-128-ctr")
	if err != nil {
		return nil, errors.New("failed to get cipher: " + err.Error())
	}

	derivedKey := eciesKDF(ecdhKey, 32)
	key_e := derivedKey[:16]
	key_m := derivedKey[16:]
	key_m_tmp := sha256.Sum256(key_m)
	key_m = key_m_tmp[:]

	iv := make([]byte, cipher.IVSize())
	_, err = rand.Read(iv)
	if err != nil {
		return nil, errors.New("failed to get random bytes: " + err.Error())
	}

	ctx, err := elliptic.NewEncryptionCipherCtx(cipher, key_e, iv)
	if err != nil {
		return nil, errors.New("failed to create cipher ctx: " + err.Error())
	}

	encData, err := ctx.Encrypt(data)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	b.Write([]byte{0x04})
	b.Write(key.PublicKey.X)
	b.Write(key.PublicKey.Y)
	b.Write(iv)
	b.Write(encData)

	hm := hmac.New(sha256.New, key_m)
	hm.Write(b.Bytes()[1+64:])
	mac := hm.Sum(nil)

	b.Write(mac)

	return b.Bytes(), nil
}

func DecryptPython(key *elliptic.PrivateKey, raw []byte) ([]byte, error) {
	cipher, err := elliptic.GetCipherByName("aes-128-ctr")
	if err != nil {
		return nil, errors.New("failed to get cipher: " + err.Error())
	}

	b := bytes.NewReader(raw)

	header4 := make([]byte, 1)
	_, err = b.Read(header4)
	if err != nil || header4[0] != 4 {
		return nil, errors.New("failed to read header")
	}

	pubkey := new(elliptic.PublicKey)
	pubkey.Curve = elliptic.Secp256k1
	pubkey.X = make([]byte, 32)
	pubkey.Y = make([]byte, 32)
	_, err = b.Read(pubkey.X)
	if err != nil {
		return nil, err
	}
	_, err = b.Read(pubkey.Y)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, cipher.IVSize())
	_, err = b.Read(iv)
	if err != nil {
		return nil, errors.New("failed to read iv")
	}

	ciphertext := make([]byte, b.Len()-sha256.Size)
	_, err = b.Read(ciphertext)
	if err != nil {
		return nil, errors.New("failed to read ciphertext")
	}

	messageMAC := make([]byte, sha256.Size)
	_, err = b.Read(messageMAC)
	if err != nil {
		return nil, errors.New("failed to read mac")
	}

	ecdhKey, err := key.GetRawECDHKey(pubkey, 32)
	if err != nil {
		return nil, errors.New("failed to get ECDH key: " + err.Error())
	}
	derivedKey := eciesKDF(ecdhKey, 32)
	key_e := derivedKey[:16]
	key_m := derivedKey[16:]
	key_m_tmp := sha256.Sum256(key_m)
	key_m = key_m_tmp[:]

	hm := hmac.New(sha256.New, key_m)
	hm.Write(raw[1+64 : len(raw)-32])
	expectedMAC := hm.Sum(nil)

	if !hmac.Equal(expectedMAC, messageMAC) {
		return nil, elliptic.InvalidMACError
	}

	ctx, err := elliptic.NewDecryptionCipherCtx(cipher, key_e, iv)
	if err != nil {
		return nil, errors.New("failed to create cipher ctx: " + err.Error())
	}

	data, err := ctx.Decrypt(ciphertext)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func eciesKDF(keyMaterial []byte, keyLen int) []byte {
	s1 := make([]byte, 0)
	key := make([]byte, 0)
	hashBlocksize := 64
	reps := ((keyLen + 7) * 8) / (hashBlocksize * 8)
	counter := 0
	buf := make([]byte, 4)
	for counter <= reps {
		counter += 1
		ctx := sha256.New()
		binary.BigEndian.PutUint32(buf, uint32(counter))
		ctx.Write(buf)
		ctx.Write(keyMaterial)
		ctx.Write(s1)
		key = append(key, ctx.Sum(nil)...)
	}
	return key[:keyLen]
}
