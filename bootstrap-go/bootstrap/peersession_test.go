package bootstrap

import (
	"encoding/hex"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ishbir/elliptic"
)

const (
	TEST_NAME     = "bootstrap-unittest"
	TEST_PROTO_ID = 1337
)

type TestAddress struct {
}

func (a *TestAddress) Network() string {
	return "test-network"
}

func (a *TestAddress) String() string {
	return "test-addr:test-port"
}

type TestConn struct {
	net.Conn
}

func (c *TestConn) RemoteAddr() net.Addr {
	return &TestAddress{}
}

type AddPeerCall struct {
	Id   string
	Peer Peer
}

type GetPeersCall struct {
	Id string
}

type TestPeerKeeper struct {
	AddPeerCalls  []AddPeerCall
	GetPeersCalls []GetPeersCall
}

func NewTestPeerKeeper() *TestPeerKeeper {
	return &TestPeerKeeper{
		AddPeerCalls:  make([]AddPeerCall, 0),
		GetPeersCalls: make([]GetPeersCall, 0),
	}
}

func (pk *TestPeerKeeper) AddPeer(id string, peer Peer) {
	pk.AddPeerCalls = append(pk.AddPeerCalls, AddPeerCall{id, peer})
}

func (pk *TestPeerKeeper) GetPeers(id string) []Peer {
	pk.GetPeersCalls = append(pk.GetPeersCalls, GetPeersCall{id})
	return nil
}

func getService(t *testing.T, pk PeerKeeper) *Service {
	privKey, err := elliptic.GeneratePrivateKey(elliptic.Secp256k1)
	if err != nil {
		t.Fatal("Error while generating private key", err)
	}

	config := &Config{
		Name:         TEST_NAME,
		Id:           "deadbeef",
		Port:         44444,
		PrvAddr:      "prvAddr",
		PubAddr:      "pubAddr",
		PrvAddresses: nil,
		NatType:      "nat type",
		PeerNum:      100,
		ProtocolId:   TEST_PROTO_ID,
	}

	return NewService(config, privKey, pk)
}

func testPeerSessionImpl(t *testing.T, handleCh chan error) {
	const (
		RAND_VAL  = 0.1337
		CLIENT_ID = "client-id"
	)
	rand.Seed(42)
	privKey, err := elliptic.GeneratePrivateKey(elliptic.Secp256k1)
	if err != nil {
		t.Fatal("Error while generating private key", err)
	}
	pubKeyHex := hex.EncodeToString(privKey.PublicKey.X) + hex.EncodeToString(privKey.PublicKey.Y)

	pk := &TestPeerKeeper{}
	service := getService(t, pk)
	conn, psConn := net.Pipe()
	ps := NewPeerSession(service, &TestConn{Conn: psConn})
	go func() {
		handleCh <- ps.handle()
	}()

	signFunc := func(msg Message) {
		sig, _ := secp256k1.Sign(GetShortHashSha(msg), privKey.Key)
		msg.GetBaseMessage().Sig = sig
	}
	encryptFunc := func(data []byte) ([]byte, error) {
		return EncryptPython(privKey, data, &service.privKey.PublicKey)
	}
	decryptFunc := func(data []byte) ([]byte, error) {
		return DecryptPython(privKey, data)
	}

	msg, err := receiveMessage(conn, nil)
	if err != nil {
		t.Fatal(err)
	}
	serverHello := msg.(*MessageHello)
	if serverHello.NodeName != TEST_NAME {
		t.Error("Wrong bootstrap node name:", serverHello.NodeName)
	}

	hello := &MessageHello{
		RandVal:     RAND_VAL,
		ClientKeyId: CLIENT_ID,
		NodeInfo: &Node{
			Key: pubKeyHex,
		},
		ProtoId: TEST_PROTO_ID,
	}
	err = sendMessage(conn, hello, encryptFunc, signFunc)
	if err != nil {
		t.Fatal(err)
	}

	randVal := &MessageRandVal{
		RandVal: serverHello.RandVal,
	}
	err = sendMessage(conn, randVal, encryptFunc, signFunc)
	if err != nil {
		t.Fatal(err)
	}

	msg, err = receiveMessage(conn, decryptFunc)
	if err != nil {
		t.Fatal(err)
	}
	serverRandVal := msg.(*MessageRandVal)
	if serverRandVal.RandVal != RAND_VAL {
		t.Fatal("Wrong RandVal", serverRandVal.RandVal)
	}

	msg, err = receiveMessage(conn, decryptFunc)
	if err != nil {
		t.Fatal(err)
	}
	serverPeers := msg.(*MessagePeers)
	if len(serverPeers.Peers) != 0 {
		t.Errorf("Expected empty list of peers, got %+v", serverPeers.Peers)
	}

	if len(pk.GetPeersCalls) != 1 {
		t.Error("GetPeers should be called once, was called:", len(pk.GetPeersCalls))
	}
	if pk.GetPeersCalls[0].Id != CLIENT_ID {
		t.Error("GetPeers was called with wrong Id:", pk.GetPeersCalls[0].Id)
	}
	if len(pk.AddPeerCalls) != 1 {
		t.Error("AddPeer should be called once, was called:", len(pk.AddPeerCalls))
	}
	if pk.AddPeerCalls[0].Id != CLIENT_ID {
		t.Error("AddPeer was called with wrong Id:", pk.AddPeerCalls[0].Id)
	}
}

func TestPeerSession(t *testing.T) {
	testCh := make(chan bool)
	handleCh := make(chan error)
	go func() {
		testPeerSessionImpl(t, handleCh)
		close(testCh)
	}()

	select {
	case <-testCh:
	case err := <-handleCh:
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("Test timed out")
	}
}
