package peer

import (
	"context"
	"io"

	"github.com/libp2p/go-libp2p-core/host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

type SidePeer interface {
	Host() host.Host
	Router() *dht.IpfsDHT
	io.Closer
}

type peer struct {
	host host.Host
	dht  *dht.IpfsDHT
}

func New(ctx context.Context, opts ...PeerOption) (SidePeer, error) {
	cfg, cfgErr := applyOptions(opts...)
	if cfgErr != nil {
		return nil, cfgErr
	}
	h, idht, err := makePeer(ctx, cfg.toListenAddr(), cfg.bootstrapPeer)
	if err != nil {
		return nil, err
	}

	ret := &peer{
		host: h,
		dht:  idht,
	}

	printListenAddr(ctx, h)
	return ret, nil
}

func (p *peer) Host() host.Host {
	return p.host
}

func (p *peer) Router() *dht.IpfsDHT {
	return p.dht
}

func (p *peer) Close() error {
	if err := p.dht.Close(); err != nil {
		return err
	}
	if err := p.host.Close(); err != nil {
		return err
	}
	return nil
}
