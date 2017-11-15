package bootstrap

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ishbir/elliptic"
	"golang.org/x/crypto/sha3"
)

type PeerSession struct {
	service *Service
	conn    net.Conn
	pubKey  *elliptic.PublicKey
	peer    Peer
	id      string
}

func NewPeerSession(service *Service, conn net.Conn) *PeerSession {
	return &PeerSession{
		service: service,
		conn:    conn,
	}
}

func (session *PeerSession) Close() {
	session.conn.Close()
}

func (session *PeerSession) sendDisconnect(reason DisconnectReason) error {
	return session.sendMessage(&MessageDisconnect{Reason: reason})
}

func (session *PeerSession) performHandshake() error {
	conn := session.conn
	service := session.service

	myHello := service.genHello()
	err := session.sendMessage(myHello)
	if err != nil {
		return err
	}
	msg, err := session.receiveMessage()
	if err != nil {
		return err
	}
	if msg.GetType() == MSG_DISCONNECT_TYPE {
		if disconnectMsg, ok := msg.(*MessageDisconnect); ok {
			return fmt.Errorf("peer disconnected, reason: %v", disconnectMsg.Reason)
		}
		return fmt.Errorf("wrong message type, was expecting Disconnect")
	}
	if msg.GetType() != MSG_HELLO_TYPE {
		return fmt.Errorf("unexpected msg type %d, was expecting Hello", msg.GetType())
	}

	helloMsg := msg.(*MessageHello)

	if helloMsg.ProtoId != session.service.config.ProtocolId {
		if err := session.sendDisconnect(DISCONNECT_PROTOCOL_VERSION); err != nil {
			return err
		}
		return fmt.Errorf("not matching protocol ID, remote %v, local %v", helloMsg.ProtoId, session.service.config.ProtocolId)
	}

	pubKeyBytes, err := hex.DecodeString(helloMsg.NodeInfo.Key)
	if err != nil {
		return fmt.Errorf("couldn't decode remote public key: %v", err)
	}
	session.pubKey, err = elliptic.PublicKeyFromUncompressedBytes(
		elliptic.Secp256k1,
		append([]byte{0x04}, pubKeyBytes...))
	if err != nil {
		return fmt.Errorf("couldn't create remote public key: %v", err)
	}

	msg, err = session.receiveMessage()
	if err != nil {
		return err
	}
	if msg.GetType() == MSG_DISCONNECT_TYPE {
		if disconnectMsg, ok := msg.(*MessageDisconnect); ok {
			return fmt.Errorf("peer disconnected, reason: %v", disconnectMsg.Reason)
		}
		return fmt.Errorf("wrong message type, was expecting Disconnect")
	}
	if msg.GetType() != MSG_RAND_VAL_TYPE {
		return fmt.Errorf("unexpected msg type %d, was expecting RandVal", msg.GetType())
	}
	randValMsg := msg.(*MessageRandVal)
	if randValMsg.RandVal != myHello.RandVal {
		return fmt.Errorf("incorrect RandVal value")
	}

	signed, err := session.verifySign(randValMsg)
	if !signed || err != nil {
		if err := session.sendDisconnect(DISCONNECT_UNVERIFIED); err != nil {
			return err
		}
		return fmt.Errorf("RandVal message not signed correctly")
	}

	myRandValMsg := MessageRandVal{RandVal: helloMsg.RandVal}
	err = session.sendMessage(&myRandValMsg)
	if err != nil {
		return err
	}

	addr, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		return err
	}

	session.peer = Peer{
		Address:  addr,
		Port:     helloMsg.Port,
		Node:     helloMsg.NodeInfo,
		NodeName: helloMsg.NodeName,
	}
	session.id = helloMsg.ClientKeyId
	return nil
}

func (session *PeerSession) handle() error {
	fmt.Println("Peer connection from", session.conn.RemoteAddr())
	err := session.performHandshake()
	if err != nil {
		return err
	}

	pk := session.service.peerKeeper
	peers := pk.GetPeers(session.id)
	peersMsg := &MessagePeers{
		Peers: make([]interface{}, len(peers)),
	}
	for idx, p := range peers {
		peersMsg.Peers[idx] = p
	}
	err = session.sendMessage(peersMsg)
	if err != nil {
		return err
	}
	pk.AddPeer(session.id, session.peer)

	disconnectMsg := &MessageDisconnect{
		Reason: DISCONNECT_BOOTSTRAP,
	}
	err = session.sendMessage(disconnectMsg)
	if err != nil {
		return err
	}

	return nil
}

func (session *PeerSession) receiveMessage() (Message, error) {
	return receiveMessage(session.conn, session.decrypt)
}

func (session *PeerSession) sendMessage(msg Message) error {
	return sendMessage(
		session.conn,
		msg,
		func(data []byte) ([]byte, error) {
			return session.encrypt(data)
		},
		session.sign)
}

func (session *PeerSession) decrypt(data []byte) ([]byte, error) {
	res, err := DecryptPython(session.service.privKey, data)
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt message: %v", err)
	}
	return res, nil
}

func (session *PeerSession) encrypt(data []byte) ([]byte, error) {
	res, err := EncryptPython(session.service.privKey, data, session.pubKey)
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt message: %v", err)
	}
	return res, nil
}

func GetShortHashSha(msg Message) []byte {
	data := msg.GetShortHash()
	sha := sha3.New256()
	sha.Write(data)
	return sha.Sum(nil)
}

func (session *PeerSession) sign(msg Message) {
	sig, _ := secp256k1.Sign(GetShortHashSha(msg), session.service.privKey.Key)
	msg.GetBaseMessage().Sig = sig
}

func (session *PeerSession) verifySign(msg Message) (bool, error) {
	keyBytes := []byte{0x04}
	keyBytes = append(keyBytes, session.pubKey.X...)
	keyBytes = append(keyBytes, session.pubKey.Y...)
	recoveredKey, err := secp256k1.RecoverPubkey(GetShortHashSha(msg), msg.GetBaseMessage().Sig)
	if err != nil {
		return false, fmt.Errorf("unable to recover public key: %v", err)
	}
	return bytes.Equal(recoveredKey, keyBytes), nil
}
