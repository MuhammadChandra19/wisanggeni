package main

import (
	"context"
	"log"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	mdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

const rendezvous = "wisanggeni" // mDNS service name

// discoveryNotifee is called every time mDNS finds a peer.
type discoveryNotifee struct{ h host.Host }

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	log.Printf("🔍 discovered %s\n", pi.ID)
	if err := n.h.Connect(context.Background(), pi); err != nil {
		log.Printf("❌ connect failed: %v", err)
	} else {
		log.Printf("✅ connected to %s", pi.ID)
	}
}

// setupMDNS enables zero‑config peer discovery on the local network.
func setupMDNS(h host.Host) error {
	s := mdns.NewMdnsService(h, rendezvous, &discoveryNotifee{h})
	return s.Start()
}

func main() {
	// 1. Start the libp2p node
	h, err := libp2p.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("👤 Peer ID  : %s", h.ID())
	for _, a := range h.Addrs() {
		log.Printf("📡 Address  : %s/p2p/%s", a, h.ID())
	}

	// 2. Enable mDNS discovery
	if err := setupMDNS(h); err != nil {
		log.Fatalf("mDNS error: %v", err)
	}

	// 3. Periodically show how many peers we’re connected to
	for {
		peers := h.Network().Peers()
		log.Printf("🌐 connected peers: %d", len(peers))
		for _, p := range peers {
			log.Printf("   • %s", p)
		}
		time.Sleep(5 * time.Second)
	}
}
