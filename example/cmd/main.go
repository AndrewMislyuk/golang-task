package main

import (
	"context"
	"log"
	"sync"

	"github.com/AndrewMislyuk/golang-task/imitation"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	imit, err := imitation.New()
	if err != nil {
		log.Fatal(err)
	}

	go imit.StopAll(cancel)

	wg.Add(3)
	go imit.Sender(ctx, &wg, "Hello")

	go imit.Sender(ctx, &wg, "Hi") // you can add second sender if you want

	go imit.Receiver(ctx, &wg)

	wg.Wait()
}
