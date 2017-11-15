package bootstrap

import (
	"sync"
)

type Peer struct {
	Address  string `cbor:"address"`
	Port     uint64 `cbor:"port"`
	Node     *Node  `cbor:"node"`
	NodeName string `cbor:"node_name"`
}

// Implementations should be thread safe.
type PeerKeeper interface {
	AddPeer(id string, peer Peer)
	GetPeers(id string) []Peer
}

type RandomizedPeerKeeper struct {
	peers   map[string]Peer
	peerNum int
	mutex   sync.Mutex
}

func NewRandomizedPeerKeeper(peerNum int) *RandomizedPeerKeeper {
	return &RandomizedPeerKeeper{
		peers:   make(map[string]Peer),
		peerNum: peerNum,
		mutex:   sync.Mutex{},
	}
}

func (pk *RandomizedPeerKeeper) AddPeer(id string, peer Peer) {
	pk.mutex.Lock()
	defer pk.mutex.Unlock()
	if _, ok := pk.peers[id]; ok {
		return
	}
	if len(pk.peers) >= pk.peerNum {
		// remove a random peer and since map iteration order is random
		// we can remove the first peer we encounter
		for id, _ := range pk.peers {
			delete(pk.peers, id)
			break
		}
	}
	pk.peers[id] = peer
}

func (pk *RandomizedPeerKeeper) GetPeers(peerId string) []Peer {
	pk.mutex.Lock()
	defer pk.mutex.Unlock()
	peers := make([]Peer, 0, len(pk.peers))
	for id, p := range pk.peers {
		if id != peerId {
			peers = append(peers, p)
		}
	}
	return peers
}
