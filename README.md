# 📄 Log Processor

**Log Processor** is a Go-based application designed to consume logs from Kafka, process them concurrently, and store them in Elasticsearch for further querying and analysis. The project exposes a REST API to allow retrieval of logs by various filters (e.g., by level or ID) and demonstrates a robust architecture for building scalable data processing systems in Go.

## 🌐 Learn More About Go for OOP Developers

Check out my dedicated web guide here: [**🧠 Go Mindset for OOP Developers**](https://rodrigogmartins.github.io/log-processor)

This site breaks down concepts, idiomatic patterns, and migration tips in a way that's easy to grasp for developers coming from Java, C#, or other OOP backgrounds.

## ⭐ Support This Project

If you found this project useful or learned something from it, consider giving it a ⭐ on GitHub!  

Your support helps keep this project alive and encourages more improvements. Thank you! 🙏

## 📘 Table of Contents

- [**1. Use Case: Centralized Log Processing System**](#1-use-case-centralized-log-processing-system)
- [**2. 🎯 Objective**](#2--objective)
- [**3. 🧠 Techniques and Learnings Applied**](#3--techniques-and-learnings-applied)
- [**4. 🚀 Running the Project**](#4--running-the-project)
- [**5. 📚 References and Best Practices Adopted**](#5--references-and-best-practices-adopted)
- [**6. 📁 Project Structure**](#6--project-structure)
- [**7. 7. 📝 Notes**](#7--notes)

## 1. Use Case: Centralized Log Processing System

A mid-size SaaS company called AcmeCloud, which operates multiple microservices written in different languages (Go, Node.js, Python) across various containers and servers. Each service generates application logs — from API requests, background jobs, and error events.

Over time, the team faces several challenges:

- 🧩 Logs are scattered across services and hard to correlate.

- 🧱 No central place to search or analyze logs.

- ⏱️ Debugging incidents takes too long because logs need to be manually aggregated.

- 📉 Monitoring system behavior and error trends is almost impossible.

To solve this, AcmeCloud decides to implement a Centralized Log Processing Pipeline using Kafka, Go, and Elasticsearch:

1. Producers (the microservices) send log events to a Kafka topic (logs).

2. The Log Processor (this Go application) consumes messages from Kafka in real time.

3. It indexes the logs in Elasticsearch, making them searchable and structured.

4. The built-in REST API allows developers to query logs easily:

    - By ID 🔑 (to track a specific request)

    - By log level 🏷️ (e.g. ERROR, INFO, WARN)

    - By time range or pattern (future extension)

With this pipeline in place, AcmeCloud’s DevOps team can:

- Quickly debug issues by searching logs in Elasticsearch/Kibana.

- Monitor trends and detect anomalies.

- Build dashboards for observability.

- Trigger alerts when certain log levels (e.g. ERROR) spike.

This use case illustrates how Log Processor helps bridge the gap between raw logs and actionable insights, turning distributed chaos into structured observability.

## 2. 🎯 Objective

The main goals of this project are:

- 📝 Build a reliable pipeline that consumes messages from Kafka and persists them in Elasticsearch.

- 🌐 Expose a REST API that allows querying and filtering logs efficiently.

- 💡 Provide a practical example of applying Go best practices in a real-world scenario.

- 📦 Create a template project that can be used as a base for future Go projects with similar requirements.

## 3. 🧠 Techniques and Learnings Applied

This project was designed to explore and consolidate multiple Go concepts and best practices:

- 🗂️ **Project structure and modularization**  
  Separation of concerns using `internal` packages (`service`, `kafka`, `api/handlers`) and a clear `cmd/main.go`.
  
- ⚡ **Concurrency in Go**  
  Processing Kafka messages concurrently using goroutines and worker pools.
  
- 🔌 **Interfaces and dependency injection**  
  Interfaces are used for Kafka consumers and Elasticsearch clients to allow easy testing and mocking.

- ✅ **Unit testing and mocks**  
  Tests utilize mocks to simulate Kafka and Elasticsearch interactions. Assert libraries were used to simplify test validations.

- ⚙️ **Configuration management**  
  Centralized configuration using `config.go` and `.env` files.

- 🛑 **Graceful shutdown**  
  Application correctly handles shutdown signals, ensuring ongoing processes finish cleanly.

- 🐳 **Dockerized dependencies**  
  Kafka, Elasticsearch, and Kafdrop are all containerized for easy local setup.

## 4. 🚀 Running the Project

1. **Clone the repository**

    ```bash
      git clone git@github.com:rodrigogmartins/log-processor.git
      cd log-processor
    ```

2. Start Docker dependencies

    ```bash
      docker-compose up -d
    ```

    This will start Kafka, Elasticsearch, and Kafdrop (Kafka UI).

3. Run the application

    ```bash
      go run cmd
    ```

4. Access the REST API

    `GET /logs` → list all logs 📄\
    `GET /logs/by-level?level=INFO` → list logs by level 🏷️\
    `GET /logs/{id}` → get log by ID 🔑

5. Optional: Run tests

    ```bash
      go test ./...
    ```

## 5. 📚 References and Best Practices Adopted

- 🗂️ Project structure: inspired by Standard Go Project Layout

- ⚡ Kafka consumer patterns: worker pools, partition-aware processing

- 🔌 ElasticSearch client usage: encapsulated in ElasticSearchClient, abstracted via interface

- ✅ Testing: dependency injection and mocks for Kafka and Elasticsearch

- ⚙️ Configuration: centralized via .env and config.go, following 12-factor principles

- 🛑 Graceful shutdown: clean termination of goroutines and external connections

## 6. 📁 Project Structure

```bash
log-processor/
│
├── cmd/
│   ├── producer/
│   │   └── main.go # Script to publish mock messages to kafka (local dev only)
│   └── main.go # Orchestrate the APP run
│
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   └── log_handler.go # API routes implementations
│   │   └── router.go # API routes declarations
│   │
│   ├── config/
│   │   └── config.go # Env configs handler
│   │
│   ├── db/
│   │   └── elastic_client.go # Elastic client abstraction
│   │
│   ├── kafka/
│   │   ├── kafka_processor.go # Kafka client connection
│   │   └── kafka_consumer.go # Consume messages logic
│   │
│   ├── service/
│   │   └── log_service.go # APP core logic
│   │
│   └── shutdown/
│       └── graceful.go # Handles grafecul shutdown
│
├── docker-compose.yaml
├── go.mod
├── go.sum
└── .env
```

## 7. 📝 Notes

This project is optimized for local development using Docker 🐳.

For production, additional considerations are needed for scaling Kafka consumers, secure Elasticsearch connections, and logging/monitoring.

Kafdrop 🎛️ is included for visualizing Kafka topics and messages during development.