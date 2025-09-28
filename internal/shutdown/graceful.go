package shutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Shutdownable interface {
	Close() error
}

func Graceful(ctx context.Context, shutdownables []Shutdownable, timeout time.Duration) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Graceful shutdown triggered...")
		cancel()

		done := make(chan struct{})
		go func() {
			for _, s := range shutdownables {
				if err := s.Close(); err != nil {
					log.Printf("Error closing resource: %v", err)
				}
			}
			close(done)
		}()

		select {
		case <-done:
			log.Println("All resources closed successfully")
		case <-time.After(timeout):
			log.Println("Timeout reached, forcing shutdown")
		}
	}()

	return ctx
}
