package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tuanuet/bwork"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	if len(os.Args) < 3 {
		log.Fatal("Usage: go run main.go <nodeID>")
	}

	target := os.Args[1]
	nodeID := os.Args[2]

	client, err := bwork.NewNodeClient(target, nodeID)
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Register(ctx); err != nil {
		log.Fatal(err)
	}

	client.StartPinging(ctx)
	eventChan := client.Subscribe(ctx)

	go func() {
		for rs := range eventChan {
			for _, event := range rs {
				fmt.Println(event)
			}
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	cancelFunc()

	log.Println("shutting down node")
}
