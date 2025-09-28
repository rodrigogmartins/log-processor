package main

import (
	"context"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/rodrigogmartins/log-processor/internal/api"
	"github.com/rodrigogmartins/log-processor/internal/config"
	"github.com/rodrigogmartins/log-processor/internal/db"
	"github.com/rodrigogmartins/log-processor/internal/kafka"
	"github.com/rodrigogmartins/log-processor/internal/service"
	"github.com/rodrigogmartins/log-processor/internal/shutdown"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using defaults or system env")
	}

	cfg := config.LoadConfig()
	ctx := context.Background()

	// --- Initializing services ---
	esClient, err := db.NewElasticSearchClient([]string{cfg.ElasticHost})
	if err != nil {
		log.Printf("Error trying to create ElasticSearchClient: %v", err)
	}

	logService := service.NewLogService(esClient, cfg.ElasticIndex)
	log.Println(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroupID)
	consumer := kafka.NewKafkaConsumer(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroupID)
	processor := kafka.NewProcessor(
		consumer,
		logService,
		cfg.MaxWorkers,
		cfg.MaxConsumeRetries,
		cfg.BackOffRetries,
	)

	// --- Inicializa graceful shutdown ---
	shutdownables := []shutdown.Shutdownable{consumer}
	ctx = shutdown.Graceful(ctx, shutdownables, cfg.ShutdownTimeout)

	// --- Rodando processor em goroutine ---
	go func() {
		log.Println("Starting Kafka processor")
		if err := processor.Start(ctx); err != nil {
			log.Printf("Kafka processor stopped with error: %v", err)
		}
	}()

	// --- Inicializa API ---
	router := api.NewRouter(esClient, cfg.ElasticIndex)
	server := &http.Server{
		Addr:    cfg.APIPort,
		Handler: router,
	}

	go func() {
		log.Printf("API server listening on %s", cfg.APIPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// --- Espera contexto ser cancelado ---
	<-ctx.Done()

	// --- Graceful shutdown HTTP ---
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("Application stopped gracefully")
}
