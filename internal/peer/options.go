package peer

import (
	"errors"
	"fmt"
	"strings"
)

const defaultPeerHost = "0.0.0.0"
const defaultPeerPort = 3001
const defaultPeerBootstrap = ""
const formatListenAddr = "/ip4/%s/tcp/%d"

type PeerOption func(*peerConfigOption)

type peerConfigOption struct {
	host          string
	port          int
	bootstrapPeer string
}

func (pco *peerConfigOption) toListenAddr() string {
	return fmt.Sprintf(formatListenAddr, pco.host, pco.port)
}

func defaultPeerConfigOption() *peerConfigOption {
	return &peerConfigOption{
		host:          defaultPeerHost,
		port:          defaultPeerPort,
		bootstrapPeer: defaultPeerBootstrap,
	}
}

func validate(pco *peerConfigOption) error {
	if pco.port < 0 || pco.port > 65535 {
		return errors.New("[sidepeer] peer configuration failed: port is invalid")
	}
	if len(pco.bootstrapPeer) < 1 {
		return errors.New("[sidepeer] peer configuration failed: bootstrap peer not specified")
	}
	return nil
}

func applyOptions(opts ...PeerOption) (*peerConfigOption, error) {
	cfg := defaultPeerConfigOption()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg, validate(cfg)
}

func WithPort(p int) PeerOption {
	return func(pco *peerConfigOption) {
		pco.port = p
	}
}

func WithBootstrapPeer(s string) PeerOption {
	return func(pco *peerConfigOption) {
		pco.bootstrapPeer = strings.TrimSpace(s)
	}
}
