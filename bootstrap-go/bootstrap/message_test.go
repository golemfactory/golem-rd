package bootstrap

import (
	"bytes"
	"testing"
)

func testImpl(t *testing.T, msg Message) Message {
	encryptCalled := false
	encryptFunc := func(data []byte) ([]byte, error) {
		encryptCalled = true
		return data, nil
	}

	decryptCalled := false
	decryptFunc := func(data []byte) ([]byte, error) {
		decryptCalled = true
		return data, nil
	}

	sig := make([]byte, SIG_LEN)
	for i := 0; i < SIG_LEN; i++ {
		sig[i] = byte(i)
	}
	signCalled := false
	signFunc := func(msg Message) {
		signCalled = true
		msg.GetBaseMessage().Sig = sig
	}

	serialized, err := serializeMessage(msg, encryptFunc, signFunc)
	if err != nil {
		t.Fatal(err)
	}
	if msg.ShouldEncrypt() != encryptCalled {
		t.Errorf("Expected encrypt func called: %v, got: %v", msg.ShouldEncrypt(), encryptCalled)
	}
	if !signCalled {
		t.Error("Sign function not called")
	}

	deserialized, err := deserializeMessage(serialized, decryptFunc)
	if err != nil {
		t.Fatal(err)
	}
	if msg.ShouldEncrypt() != decryptCalled {
		t.Errorf("Expected decrypt func called: %v, got: %v", msg.ShouldEncrypt(), decryptCalled)
	}

	baseMsg := deserialized.GetBaseMessage()
	if baseMsg.Header.Type != msg.GetType() {
		t.Error("Wrong message type, expected %v, got %v", msg.GetType(), baseMsg.Header.Type)
	}
	if baseMsg.Header.Encrypted != msg.ShouldEncrypt() {
		t.Errorf("Expected encrypted field to be: %v, got: %v", msg.ShouldEncrypt(), baseMsg.Header.Encrypted)
	}
	if !bytes.Equal(sig, baseMsg.Sig) {
		t.Errorf("Wrong signature, expected %v, got %v", sig, baseMsg.Sig)
	}
	return deserialized
}

func TestSerializeationEncrypted(t *testing.T) {
	const RAND_VAL = 0.1337
	msg := &MessageRandVal{
		RandVal: RAND_VAL,
	}
	if !msg.ShouldEncrypt() {
		t.Fatal("Tested message should be encryptable")
	}
	deserialized := testImpl(t, msg)

	castedMsg, ok := deserialized.(*MessageRandVal)
	if !ok {
		t.Fatal("Message should be of type MessageRandVal")
	}
	if castedMsg.RandVal != RAND_VAL {
		t.Errorf("Wrong rand val, expected %v, got %v", RAND_VAL, castedMsg.RandVal)
	}
}

func TestSerializeationNotEncrypted(t *testing.T) {
	const REASON = "Unittest"
	msg := &MessageDisconnect{
		Reason: REASON,
	}
	if msg.ShouldEncrypt() {
		t.Fatal("Tested message shouldn't be encryptable")
	}
	deserialized := testImpl(t, msg)

	castedMsg, ok := deserialized.(*MessageDisconnect)
	if !ok {
		t.Fatal("Message should be of type MessageDisconnect")
	}
	if castedMsg.Reason != REASON {
		t.Errorf("Wrong rand val, expected %v, got %v", REASON, castedMsg.Reason)
	}
}
