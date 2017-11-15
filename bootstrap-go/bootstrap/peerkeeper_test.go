package bootstrap

import (
	"testing"
)

func TestRandomizedPeerKeeper(t *testing.T) {
	pk := NewRandomizedPeerKeeper(2)
	peers := pk.GetPeers("foo")
	if len(peers) != 0 {
		t.Errorf("Expected empty list of peers, got %v", peers)
	}

	peer1 := Peer{NodeName: "peer1"}
	pk.AddPeer("peer1", peer1)
	peers = pk.GetPeers("foo")
	if len(peers) != 1 || peers[0].NodeName != "peer1" {
		t.Errorf("Expected peer1, got %v", peers)
	}

	peer2 := Peer{NodeName: "peer2"}
	pk.AddPeer("peer2", peer2)
	peers = pk.GetPeers("foo")
	if len(peers) != 2 {
		t.Errorf("Expected two peers, got %v", peers)
	}

	peers = pk.GetPeers("peer2")
	if len(peers) != 1 || peers[0].NodeName != "peer1" {
		t.Errorf("Expected peer1, got %v", peers)
	}

	peer3 := Peer{NodeName: "peer3"}
	pk.AddPeer("peer3", peer3)
	peers = pk.GetPeers("foo")
	if len(peers) != 2 {
		t.Errorf("Expected two peers, got %v", peers)
	}
	if peers[0].NodeName != "peer3" && peers[1].NodeName != "peer3" {
		t.Errorf("Expected peer3 to be in the list, got %v", peers)
	}
}
