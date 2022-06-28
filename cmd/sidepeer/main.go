package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	igrpc "github.com/igumus/sidepeer/internal/grpc"
	ipeer "github.com/igumus/sidepeer/internal/peer"
)

func main() {
	log.Println("info: starting sidepeer")
	flagGrpcPort := flag.Int("grpc-port", 4001, "GrpcServer Port Addr")
	flagPeerPort := flag.Int("peer-port", 3002, "Peer Port Addr")
	flagPeerBootstrap := flag.String("peer-bootstrap", "", "Peer Bootstrap Addr")
	flag.Parse()

	rootCtx := context.Background()

	peer, err := ipeer.New(rootCtx, ipeer.WithPort(*flagPeerPort), ipeer.WithBootstrapPeer(*flagPeerBootstrap))
	if err != nil {
		log.Fatalf("err: creating peer failed: %s\n", err.Error())
	}

	grpcServer, err := igrpc.NewGrpcServer(rootCtx, igrpc.WithPort(*flagGrpcPort))
	if err != nil {
		log.Fatalf("err: creating grpc server failed: %s\n", err.Error())
	}

	grpcServer.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-c
	log.Println("err: shutdown process started")
	grpcServer.Close()
	peer.Close()
}
