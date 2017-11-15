package bootstrap

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"

	"github.com/ishbir/elliptic"
)

type Config struct {
	Name         string
	Id           string
	Port         uint64
	PrvAddr      string
	PubAddr      string
	PrvAddresses []interface{}
	NatType      string
	PeerNum      int
	ProtocolId   uint64
}

type Service struct {
	config     *Config
	privKey    *elliptic.PrivateKey
	pubKeyHex  string
	peerKeeper PeerKeeper
}

func NewService(config *Config, privKey *elliptic.PrivateKey, pk PeerKeeper) *Service {
	pubKeyHex := hex.EncodeToString(privKey.PublicKey.X) + hex.EncodeToString(privKey.PublicKey.Y)
	return &Service{
		config:     config,
		privKey:    privKey,
		pubKeyHex:  pubKeyHex,
		peerKeeper: pk,
	}
}

func (s *Service) Listen() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		return err
	}
	defer l.Close()
	fmt.Printf("Listening on port %d\n", s.config.Port)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
		} else {
			go func() {
				ps := NewPeerSession(s, conn)
				err := ps.handle()
				ps.Close()
				if err != nil {
					fmt.Println("Peer session error:", err)
				}
			}()
		}
	}
	return nil
}

func (s *Service) genHello() *MessageHello {
	return &MessageHello{
		Port:        s.config.Port,
		NodeName:    s.config.Name,
		ClientKeyId: s.config.Id,
		NodeInfo: &Node{
			NodeName:     s.config.Name,
			Key:          s.pubKeyHex,
			PrvPort:      0,
			PubPort:      0,
			P2pPrvPort:   s.config.Port,
			P2pPubPort:   0,
			PrvAddr:      s.config.PrvAddr,
			PubAddr:      s.config.PubAddr,
			PrvAddresses: s.config.PrvAddresses,
			NatType:      s.config.NatType,
		},
		RandVal:        rand.Float64(),
		Metadata:       make(map[string]interface{}),
		SolveChallange: false,
		Challange:      nil,
		Difficulty:     0,
		ProtoId:        s.config.ProtocolId,
		ClientVer:      "0.9.0",
	}
}
