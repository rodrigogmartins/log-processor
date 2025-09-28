package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Kafka
	KafkaBrokers         []string
	KafkaTopic           string
	KafkaGroupID         string
	MaxWorkers           int
	MaxConsumeRetries    int
	BackOffRetries       time.Duration
	WorkerTimeoutSeconds int

	// Elasticsearch
	ElasticHost  string
	ElasticIndex string

	// API
	APIPort string

	// Other
	ShutdownTimeout time.Duration
}

func LoadConfig() *Config {
	maxWorkers, err := strconv.Atoi(os.Getenv("KAFKA_MAX_WORKERS"))
	if err != nil {
		maxWorkers = 10
	}

	maxConsumeRetries, err := strconv.Atoi(os.Getenv("KAFKA_MAX_CONSUME_RETRIES"))
	if err != nil {
		maxConsumeRetries = 3
	}

	backOffRetriesMs, err := strconv.Atoi(os.Getenv("KAFKA_BACKOFF_TIME_SECONDS"))
	if err != nil {
		backOffRetriesMs = 300
	}

	workerTimeout, err := strconv.Atoi(os.Getenv("WORKER_TIMEOUT_SECONDS"))
	if err != nil {
		workerTimeout = 1
	}

	timeoutStr := os.Getenv("SHUTDOWN_TIMEOUT")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		log.Printf("Invalid SHUTDOWN_TIMEOUT, using default 5s: %v", err)
		timeout = 5 * time.Second
	}

	return &Config{
		KafkaBrokers:         []string{os.Getenv("KAFKA_BROKERS")},
		KafkaTopic:           os.Getenv("KAFKA_TOPIC"),
		KafkaGroupID:         os.Getenv("KAFKA_GROUP_ID"),
		MaxWorkers:           maxWorkers,
		MaxConsumeRetries:    maxConsumeRetries,
		BackOffRetries:       (time.Duration(backOffRetriesMs) * time.Millisecond),
		WorkerTimeoutSeconds: workerTimeout,
		ElasticHost:          os.Getenv("ELASTIC_HOST"),
		ElasticIndex:         os.Getenv("ELASTIC_INDEX"),
		APIPort:              os.Getenv("API_PORT"),
		ShutdownTimeout:      timeout,
	}
}
