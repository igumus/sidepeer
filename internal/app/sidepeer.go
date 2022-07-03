package app

import (
	"context"
	"io"
	"log"

	"github.com/igumus/blockstorage"

	"github.com/igumus/blockstorage/blockpb"
	bgrpc "github.com/igumus/blockstorage/grpc"
	bpeer "github.com/igumus/blockstorage/peer"
	fsstore "github.com/igumus/go-objectstore-fs"
	igrpc "github.com/igumus/sidepeer/internal/grpc"
	ipeer "github.com/igumus/sidepeer/internal/peer"
)

type App interface {
	Start() error
	io.Closer
}

type app struct {
	peer ipeer.SidePeer
	grpc igrpc.GrpcServer
}

func NewApp(ctx context.Context, grpcPort, peerPort int, peerBootstrap, dataDir, bucketDir string) (App, error) {
	ret := &app{}

	if err := ret.makeSidePeer(ctx, peerPort, peerBootstrap); err != nil {
		return nil, err
	}

	if err := ret.makeGrpcServer(ctx, grpcPort); err != nil {
		return nil, err
	}

	if err := ret.makeBlockStorage(ctx, dataDir, bucketDir); err != nil {
		return nil, err
	}

	return ret, nil
}

func (a *app) makeSidePeer(ctx context.Context, peerPort int, peerBootstrap string) error {
	peer, err := ipeer.New(ctx, ipeer.WithPort(peerPort), ipeer.WithBootstrapPeer(peerBootstrap))
	if err != nil {
		return err
	}
	a.peer = peer
	return nil
}

func (a *app) makeGrpcServer(ctx context.Context, grpcPort int) error {
	grpcServer, err := igrpc.New(ctx, grpcPort)
	if err != nil {
		return err
	}

	a.grpc = grpcServer
	return nil
}

func (a *app) makeBlockStorage(ctx context.Context, dirData, dirBucket string) error {
	dataDirOption := fsstore.WithDataDir(dirData)

	store, storeErr := fsstore.NewFileSystemObjectStore(dataDirOption, fsstore.WithBucket(dirBucket))
	if storeErr != nil {
		return storeErr
	}

	tempStore, tempStoreErr := fsstore.NewFileSystemObjectStore(dataDirOption, fsstore.WithBucket("temporary"))
	if tempStoreErr != nil {
		return tempStoreErr
	}

	blockstoragePeer, err := bpeer.NewBlockStoragePeer(ctx, bpeer.WithHost(a.peer.Host()), bpeer.WithContentRouter(a.peer.Router()), bpeer.WithTempStore(tempStore))
	if err != nil {
		return err
	}

	blockstorageService, err := blockstorage.NewBlockStorage(ctx, blockstorage.WithLocalStore(store), blockstorage.WithPeer(blockstoragePeer))
	if err != nil {
		return err
	}

	blockstorageEndpoint, err := bgrpc.NewBlockStorageServiceEndpoint(ctx, blockstorageService)
	if err != nil {
		return err
	}

	blockpb.RegisterBlockStorageGrpcServiceServer(a.grpc.Server(), blockstorageEndpoint)
	return nil
}

func (a *app) Start() error {
	return a.grpc.Start()
}

func (a *app) Close() error {
	log.Println("info: stopping the application")
	if err := a.grpc.Close(); err != nil {
		return err
	}

	if err := a.peer.Close(); err != nil {
		return err
	}

	return nil
}
