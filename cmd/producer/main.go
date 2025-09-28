package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/rodrigogmartins/log-processor/internal/config"
	"github.com/segmentio/kafka-go"
)

type LogMessage struct {
	ID        string `json:"id"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using defaults or system env")
	}
	cfg := config.LoadConfig()
	kafkaBroker := cfg.KafkaBrokers[0]

	err := ensureTopicExists(kafkaBroker, cfg.KafkaTopic)
	if err != nil {
		log.Fatalf("Error trying to check the topic %v", err)
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  cfg.KafkaBrokers,
		Topic:    cfg.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	for i := 1; i <= 5; i++ {
		msg := LogMessage{
			ID:        fmt.Sprintf("%v-%v", time.Now().Format("20060102150405"), i),
			Level:     "INFO",
			Message:   fmt.Sprintf("Test message number %v", i),
			Timestamp: time.Now().Format(time.RFC3339),
		}

		data, _ := json.Marshal(msg)
		err := writer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(msg.ID),
			Value: data,
		})
		if err != nil {
			log.Fatalf("Error trying to publish message: %v", err)
		}

		log.Printf("Message sent: %+v\n", msg)
	}

	log.Println("Messages sent with success")
}

func ensureTopicExists(broker, topic string) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return err
	}
	defer conn.Close()

	controllerConn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	return controllerConn.CreateTopics(topicConfigs...)
}
