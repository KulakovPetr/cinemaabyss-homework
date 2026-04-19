package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

var (
	producer sarama.SyncProducer
	consumer sarama.Consumer
)

type MovieEvent struct {
	MovieID    int      `json:"movie_id"`
	Title      string   `json:"title"`
	Action     string   `json:"action"`
	UserID     *int     `json:"user_id,omitempty"`
	Rating     *float64 `json:"rating,omitempty"`
	Genres     []string `json:"genres,omitempty"`
	Description *string `json:"description,omitempty"`
}

type UserEvent struct {
	UserID    int    `json:"user_id"`
	Username  *string `json:"username,omitempty"`
	Email     *string `json:"email,omitempty"`
	Action    string `json:"action"`
	Timestamp string `json:"timestamp"`
}

type PaymentEvent struct {
	PaymentID  int     `json:"payment_id"`
	UserID     int     `json:"user_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	Timestamp  string  `json:"timestamp"`
	MethodType *string `json:"method_type,omitempty"`
}

type EventResponse struct {
	Status   string      `json:"status"`
	Partition int32      `json:"partition"`
	Offset   int64       `json:"offset"`
	Event    interface{} `json:"event"`
}

func main() {
	// Read environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "kafka:9092"
	}

	// Initialize Kafka producer
	var err error
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err = sarama.NewSyncProducer([]string{kafkaBrokers}, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Initialize Kafka consumer
	consumer, err = sarama.NewConsumer([]string{kafkaBrokers}, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Start consuming events in background
	go startConsumers()

	// Set up HTTP routes
	http.HandleFunc("/api/events/health", healthHandler)
	http.HandleFunc("/api/events/movie", handleMovieEvent)
	http.HandleFunc("/api/events/user", handleUserEvent)
	http.HandleFunc("/api/events/payment", handlePaymentEvent)

	log.Printf("Starting events service on port %s", port)
	log.Printf("Kafka brokers: %s", kafkaBrokers)

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down...")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"status": true})
}

func handleMovieEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var event MovieEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Create Kafka message
	eventJSON, _ := json.Marshal(event)
	msg := &sarama.ProducerMessage{
		Topic: "movie-events",
		Value: sarama.StringEncoder(eventJSON),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		http.Error(w, fmt.Sprintf("Failed to publish event: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Published movie event to Kafka: partition=%d, offset=%d, event=%+v", partition, offset, event)

	response := EventResponse{
		Status:   "success",
		Partition: partition,
		Offset:   offset,
		Event:    event,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func handleUserEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var event UserEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Create Kafka message
	eventJSON, _ := json.Marshal(event)
	msg := &sarama.ProducerMessage{
		Topic: "user-events",
		Value: sarama.StringEncoder(eventJSON),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		http.Error(w, fmt.Sprintf("Failed to publish event: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Published user event to Kafka: partition=%d, offset=%d, event=%+v", partition, offset, event)

	response := EventResponse{
		Status:   "success",
		Partition: partition,
		Offset:   offset,
		Event:    event,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func handlePaymentEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var event PaymentEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Create Kafka message
	eventJSON, _ := json.Marshal(event)
	msg := &sarama.ProducerMessage{
		Topic: "payment-events",
		Value: sarama.StringEncoder(eventJSON),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		http.Error(w, fmt.Sprintf("Failed to publish event: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Published payment event to Kafka: partition=%d, offset=%d, event=%+v", partition, offset, event)

	response := EventResponse{
		Status:   "success",
		Partition: partition,
		Offset:   offset,
		Event:    event,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func startConsumers() {
	topics := []string{"movie-events", "user-events", "payment-events"}

	for _, topic := range topics {
		go consumeTopic(topic)
	}
}

func consumeTopic(topic string) {
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Printf("Failed to start consumer for topic %s: %v", topic, err)
		return
	}
	defer partitionConsumer.Close()

	log.Printf("Started consuming from topic: %s", topic)

	for {
		select {
		case message := <-partitionConsumer.Messages():
			if message != nil {
				log.Printf("Consumed event from topic %s: partition=%d, offset=%d, value=%s",
					topic, message.Partition, message.Offset, string(message.Value))
			}
		case err := <-partitionConsumer.Errors():
			if err != nil {
				log.Printf("Error consuming from topic %s: %v", topic, err)
			}
		}
	}
}
