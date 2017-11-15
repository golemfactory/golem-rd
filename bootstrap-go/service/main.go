package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"golem_rd/bootstrap-go/bootstrap"
	"math/rand"
	"net"
	"time"

	"github.com/ccding/go-stun/stun"
	"github.com/ishbir/elliptic"
)

const (
	PORT            = 40102
	PEER_NUM        = 100
	NAME            = "Go Bootstrap"
	PROTO_ID uint64 = 18
)

func main() {
	var port uint64
	var peerNum int
	var name string
	var protocolId uint64
	flag.Uint64Var(&port, "port", PORT, "Port to listen to")
	flag.IntVar(&peerNum, "peer-num", PEER_NUM, "Number of peers to send")
	flag.StringVar(&name, "name", NAME, "Name of the node")
	flag.Uint64Var(&protocolId, "protocol-id", PROTO_ID, "Version of the P2P procotol")
	flag.Parse()

	var err error
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error getting network interfaces:", err)
		return
	}
	prvAddresses := make([]interface{}, 0)
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				prvAddresses = append(prvAddresses, ipnet.IP.String())
			}
		}
	}

	nat, host, err := stun.NewClient().Discover()
	if err != nil {
		fmt.Println("Error discovering STUN details:", err)
		return
	}

	rand.Seed(time.Now().UTC().UnixNano())
	privKey, err := elliptic.GeneratePrivateKey(elliptic.Secp256k1)
	if err != nil {
		fmt.Println("Error while generating private key", err)
		return
	}
	pubKeyHex := hex.EncodeToString(privKey.PublicKey.X) + hex.EncodeToString(privKey.PublicKey.Y)

	config := &bootstrap.Config{
		Name:         name,
		Id:           pubKeyHex,
		Port:         port,
		PrvAddr:      prvAddresses[0].(string),
		PubAddr:      host.IP(),
		PrvAddresses: prvAddresses,
		NatType:      nat.String(),
		PeerNum:      PEER_NUM,
		ProtocolId:   protocolId,
	}

	fmt.Printf("Config: %+v\n", config)

	service := bootstrap.NewService(
		config,
		privKey,
		bootstrap.NewRandomizedPeerKeeper(config.PeerNum))
	err = service.Listen()
	if err != nil {
		fmt.Println("Error during listen:", err)
	}
}
