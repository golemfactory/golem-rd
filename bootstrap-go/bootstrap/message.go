package bootstrap

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"reflect"
	"time"

	"github.com/whyrusleeping/cbor/go"
)

type Message interface {
	GetBaseMessage() *BaseMessage
	GetShortHash() []byte
	GetType() uint16
	ShouldEncrypt() bool
}

type BaseMessage struct {
	Header            Header
	Sig               []byte
	serializedPayload []byte
}

func (self *BaseMessage) GetBaseMessage() *BaseMessage {
	return self
}

func (self *BaseMessage) GetShortHash() []byte {
	data := make([]byte, 0)
	data = append(data, self.Header.serialize()...)
	data = append(data, self.serializedPayload...)
	hash := sha1.Sum(data)
	return hash[:]
}

// slot is a pair fo python field's name and value
type MessageSlot = []interface{}

// list of MessageSlots
type MessagePayload = []interface{}

func GetPayload(msg Message) MessagePayload {
	payload := MessagePayload{}
	v := reflect.ValueOf(msg).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		val := v.Field(i)
		tag := field.Tag.Get("msg_slot")
		if tag != "" {
			payload = append(payload, MessageSlot{tag, val.Interface()})
		}
	}
	return payload
}

// https://pypi.python.org/pypi/cbor2/3.0.4
// Semantics: Mark shared value
type Tag28Decoder struct{}

func (self *Tag28Decoder) GetTag() uint64 {
	return 28
}

func (self *Tag28Decoder) DecodeTarget() interface{} {
	var v interface{}
	return &v
}

func (self *Tag28Decoder) PostDecode(v interface{}) (interface{}, error) {
	return *v.(*interface{}), nil
}

const (
	HEADER_LEN = 11
	SIG_LEN    = 65
)

type Header struct {
	Type      uint16
	Timestamp uint64
	Encrypted bool
}

func (self *Header) serialize() []byte {
	res := make([]byte, HEADER_LEN)
	binary.BigEndian.PutUint16(res, self.Type)
	binary.BigEndian.PutUint64(res[2:], self.Timestamp)
	if self.Encrypted {
		res[10] = 1
	}
	return res
}

func deserializeHeader(header []byte) Header {
	typ := binary.BigEndian.Uint16(header[:2])
	timestamp := binary.BigEndian.Uint64(header[2:10])
	encrypted := header[10] == 1
	return Header{typ, timestamp, encrypted}
}

type DecryptFunc = func([]byte) ([]byte, error)

func deserializeMessage(b []byte, decrypt DecryptFunc) (Message, error) {
	payloadIdx := HEADER_LEN + SIG_LEN
	headerB := b[:HEADER_LEN]
	sigB := b[HEADER_LEN:payloadIdx]
	payloadB := b[payloadIdx:]

	header := deserializeHeader(headerB)
	var msg Message
	if header.Type == MSG_HELLO_TYPE {
		msg = &MessageHello{}
	} else if header.Type == MSG_RAND_VAL_TYPE {
		msg = &MessageRandVal{}
	} else if header.Type == MSG_DISCONNECT_TYPE {
		msg = &MessageDisconnect{}
	} else if header.Type == MSG_PEERS_TYPE {
		msg = &MessagePeers{}
	} else {
		return nil, fmt.Errorf("unsupported msg type %d", header.Type)
	}

	msg.GetBaseMessage().Header = header
	msg.GetBaseMessage().Sig = sigB
	msg.GetBaseMessage().serializedPayload = payloadB

	var err error
	if header.Encrypted {
		payloadB, err = decrypt(payloadB)
		if err != nil {
			return nil, err
		}
	}

	reader := bytes.NewReader(payloadB)
	decoder := cbor.NewDecoder(reader)
	var pyObjectDecoder PyObjectDecoder
	decoder.TagDecoders[pyObjectDecoder.GetTag()] = &pyObjectDecoder
	var tag28Decoder Tag28Decoder
	decoder.TagDecoders[tag28Decoder.GetTag()] = &tag28Decoder
	var slots MessagePayload
	err = decoder.Decode(&slots)
	if err != nil {
		return nil, err
	}
	deserializePayload(slots, msg)
	return msg, nil
}

func deserializePayload(slotsList MessagePayload, msg Message) {
	slots := make(map[string]interface{})
	for _, s := range slotsList {
		slot, ok := s.(MessageSlot)
		if !ok {
			fmt.Printf("Couldn't cast slot %+v\n", s)
			continue
		}
		slots[slot[0].(string)] = slot[1]
	}

	v := reflect.ValueOf(msg).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		val := v.Field(i)
		tag := field.Tag.Get("msg_slot")
		if tag != "" {
			if vv, ok := slots[tag]; ok && vv != nil {
				val.Set(reflect.ValueOf(vv))
			}
		}
	}
}

func serializePayload(payload MessagePayload) ([]byte, error) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	encoder := cbor.NewEncoder(writer)
	err := encoder.Encode(payload)
	if err != nil {
		return nil, err
	}
	writer.Flush()
	return b.Bytes(), nil
}

type EncryptFunc = func([]byte) ([]byte, error)
type SignFunc = func(Message)

func serializeMessage(msg Message, encrypt EncryptFunc, sign SignFunc) ([]byte, error) {
	header := &msg.GetBaseMessage().Header
	header.Type = msg.GetType()
	header.Timestamp = uint64(time.Now().UnixNano()) / uint64(time.Microsecond)
	header.Encrypted = msg.ShouldEncrypt()
	headerBytes := header.serialize()
	payloadBytes, err := serializePayload(GetPayload(msg))
	if msg.ShouldEncrypt() {
		payloadBytes, err = encrypt(payloadBytes)
		if err != nil {
			return nil, err
		}
	}
	msg.GetBaseMessage().serializedPayload = payloadBytes
	sign(msg)
	sigBytes := msg.GetBaseMessage().Sig

	if err != nil {
		return nil, err
	}
	res := make([]byte, 0)
	res = append(res, headerBytes...)
	res = append(res, sigBytes...)
	res = append(res, payloadBytes...)
	return res, nil
}

const (
	MSG_HELLO_TYPE      = 0
	MSG_RAND_VAL_TYPE   = 1
	MSG_DISCONNECT_TYPE = 2
	MSG_PEERS_TYPE      = 1004
)

type MessageHello struct {
	BaseMessage
	Port           uint64      `msg_slot:"port"`
	NodeName       string      `msg_slot:"node_name"`
	ClientKeyId    string      `msg_slot:"client_key_id"`
	NodeInfo       *Node       `msg_slot:"node_info"`
	RandVal        float64     `msg_slot:"rand_val"`
	Metadata       interface{} `msg_slot:"metadata"`
	SolveChallange bool        `msg_slot:"solve_challenge"`
	Challange      interface{} `msg_slot:"challenge"`
	Difficulty     uint64      `msg_slot:"difficulty"`
	ProtoId        uint64      `msg_slot:"proto_id"`
	ClientVer      string      `msg_slot:"client_ver"`
}

func (self *MessageHello) GetType() uint16 {
	return MSG_HELLO_TYPE
}

func (self *MessageHello) ShouldEncrypt() bool {
	return false
}

type MessageRandVal struct {
	BaseMessage
	RandVal float64 `msg_slot:"rand_val"`
}

func (self *MessageRandVal) GetType() uint16 {
	return MSG_RAND_VAL_TYPE
}

func (self *MessageRandVal) ShouldEncrypt() bool {
	return true
}

type DisconnectReason = string

const (
	DISCONNECT_PROTOCOL_VERSION DisconnectReason = "protocol_version"
	DISCONNECT_UNVERIFIED       DisconnectReason = "unverified"
	DISCONNECT_BOOTSTRAP        DisconnectReason = "bootstrap"
)

type MessageDisconnect struct {
	BaseMessage
	Reason DisconnectReason `msg_slot:"reason"`
}

func (self *MessageDisconnect) GetType() uint16 {
	return MSG_DISCONNECT_TYPE
}

func (self *MessageDisconnect) ShouldEncrypt() bool {
	return false
}

type MessagePeers struct {
	BaseMessage
	Peers []interface{} `msg_slot:"peers"`
}

func (self *MessagePeers) GetType() uint16 {
	return MSG_PEERS_TYPE
}

func (self *MessagePeers) ShouldEncrypt() bool {
	return true
}
