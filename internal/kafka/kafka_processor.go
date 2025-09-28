package kafka

import (
	"context"
	"io"
	"log"
	"sync"
	"time"

	"github.com/rodrigogmartins/log-processor/internal/service"
	"github.com/segmentio/kafka-go"
)

type KafkaReader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
	CommitMessage(msg kafka.Message) error
}

type Processor struct {
	Reader       KafkaReader
	LogService   service.LogServiceInterface
	MaxWorkers   int
	RetryMax     int
	RetryBackoff time.Duration
}

func NewProcessor(reader KafkaReader, logService service.LogServiceInterface, maxWorkers int, retryMax int, retryBackoff time.Duration) *Processor {
	return &Processor{
		Reader:       reader,
		LogService:   logService,
		MaxWorkers:   maxWorkers,
		RetryMax:     retryMax,
		RetryBackoff: retryBackoff,
	}
}

func (p *Processor) Start(ctx context.Context) error {
	sem := make(chan struct{}, p.MaxWorkers)
	var wg sync.WaitGroup

	for {
		log.Println("Kafka processor reading message")
		msg, err := p.Reader.ReadMessage(ctx)
		log.Printf("read message: %v", msg)
		if err != nil {
			if err == io.EOF {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			if ctx.Err() != nil {
				break
			}
			log.Printf("Error reading message: %v", err)
			continue
		}

		sem <- struct{}{}
		wg.Add(1)

		go func(m kafka.Message) {
			defer func() {
				<-sem
				wg.Done()
			}()

			logEntry := service.Log{
				ID:      string(m.Key),
				Message: string(m.Value),
			}

			for i := 0; i < p.RetryMax; i++ {
				opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

				err := p.LogService.Process(opCtx, logEntry)
				cancel()
				if err == nil {
					break
				}

				log.Printf("Retry %d: error processing log %s: %v", i+1, logEntry.ID, err)
				time.Sleep(p.RetryBackoff * time.Duration(i+1))
			}

			if err := p.Reader.CommitMessage(m); err != nil {
				log.Printf("Error committing message: %v", err)
			}
		}(msg)
	}

	wg.Wait()
	return p.Reader.Close()
}
