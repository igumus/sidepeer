package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	iapp "github.com/igumus/sidepeer/internal/app"
)

func main() {
	log.Println("info: starting sidepeer")

	flagPeerPort := flag.Int("peer-port", 4001, "Peer Port Addr")
	flagPeerBootstrap := flag.String("peer-bootstrap", "", "Peer Bootstrap Addr")
	flagGrpcPort := flag.Int("grpc-port", 4002, "GrpcServer Port Addr")
	flagDirData := flag.String("store-dir-data", "/tmp", "ObjectStore Data Dir")
	flagDirBucket := flag.String("store-dir-bucket", "objects", "ObjectStore Bucket Dir")
	flag.Parse()

	rootCtx := context.Background()

	app, err := iapp.NewApp(rootCtx, *flagGrpcPort, *flagPeerPort, *flagPeerBootstrap, *flagDirData, *flagDirBucket)
	if err != nil {
		log.Fatalf("err: creating application failed: %s\n", err.Error())
	}

	if err := app.Start(); err != nil {
		log.Fatalf("err: starting application failed: %s\n", err.Error())
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-c
	log.Println("err: shutdown process started")
	app.Close()
}
