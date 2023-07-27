package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func (app *Config) runConsumer() {
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

loop:
	for {
		select {
		case <-signalChan: // if get SIGTERM
			log.Println("Got SIGTERM signal, cancelling the context")
			cancel() //cancel context

		default:
			_, err := app.processSQS(ctx)

			if err != nil {
				if errors.Is(err, context.Canceled) {
					log.Printf("stop processing, context is cancelled %v", err)
					break loop
				}

				log.Fatalf("error processing SQS %v", err)
			}
		}
	}
}
