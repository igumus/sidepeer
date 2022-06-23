package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("info: starting sidepeer")
	// creating root context
	rootCtx := context.Background()
	log.Printf("info: root context created")

	log.Println("info: root context: ", rootCtx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-c
	log.Println("err: shutdown process started")
	// grpcServer.Shutdown()
}
