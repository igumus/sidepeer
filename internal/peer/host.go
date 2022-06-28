package peer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	libpeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	routedhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	ma "github.com/multiformats/go-multiaddr"
)

func printListenAddr(ctx context.Context, host host.Host) {
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", host.ID().Pretty()))
	for _, addr := range host.Addrs() {
		log.Printf("info: peer listening on : %s\n", addr.Encapsulate(hostAddr).String())
	}
}

func generateKeyPair(ctx context.Context) (crypto.PrivKey, error) {
	sk, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		return nil, err
	}
	return sk, nil
}

func convertPeers(peers []string) []libpeer.AddrInfo {
	pinfos := make([]libpeer.AddrInfo, len(peers))
	for i, addr := range peers {
		maddr := ma.StringCast(addr)
		p, err := libpeer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			log.Fatalln(err)
		}
		pinfos[i] = *p
	}
	return pinfos
}

func connectBootstrapPeer(ctx context.Context, ph host.Host, peers ...libpeer.AddrInfo) error {
	if len(peers) < 1 {
		return errors.New("not enough bootstrap peers")
	}

	errs := make(chan error, len(peers))
	var wg sync.WaitGroup
	for _, p := range peers {
		wg.Add(1)
		go func(p libpeer.AddrInfo) {
			defer wg.Done()
			ph.Peerstore().AddAddrs(p.ID, p.Addrs, peerstore.PermanentAddrTTL)
			if err := ph.Connect(ctx, p); err != nil {
				log.Printf("err: failed to connect bootstrap peer: %s, %s\n", p.ID, err.Error())
				errs <- err
				return
			}
			log.Printf("info: succeded to connect bootstrap peer: %s\n", p.ID)
		}(p)
	}
	wg.Wait()

	close(errs)
	count := 0
	var err error
	for err = range errs {
		if err != nil {
			count++
		}
	}
	if count == len(peers) {
		return fmt.Errorf("failed to bootstrap. %s", err)
	}
	return nil
}

func makeHost(ctx context.Context, listenAddr string) (host.Host, error) {
	sk, err := generateKeyPair(ctx)
	if err != nil {
		log.Printf("err: generation key pair failed: %s\n", err.Error())
		return nil, err
	}
	connmgr, err := connmgr.NewConnManager(
		100,
		400,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		return nil, err
	}

	host, err := libp2p.New(
		libp2p.Identity(sk),
		libp2p.ListenAddrStrings(listenAddr),
		libp2p.ConnectionManager(connmgr),
		libp2p.DefaultTransports,
	)
	if err != nil {
		return nil, err
	}
	return host, nil
}

func makePeer(ctx context.Context, listenAddr, bootstrapAddr string) (host.Host, *dht.IpfsDHT, error) {
	log.Printf("info: new peer with listen addr: %s\n", listenAddr)

	host, err := makeHost(ctx, listenAddr)
	if err != nil {
		return nil, nil, err
	}
	bootstrapHosts := convertPeers([]string{bootstrapAddr})

	idht, err := dht.New(ctx, host, dht.BootstrapPeersFunc(func() []libpeer.AddrInfo {
		return bootstrapHosts
	}))
	if err != nil {
		return host, nil, err
	}

	connectBootstrapPeer(ctx, host, bootstrapHosts...)

	if err := idht.Bootstrap(ctx); err != nil {
		log.Printf("warn: dht bootstrapping failed: %s\n", err.Error())
	}

	rhost := routedhost.Wrap(host, idht)
	return rhost, idht, nil
}
