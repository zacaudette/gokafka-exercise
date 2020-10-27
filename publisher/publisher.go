package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

const (
	address   = "kafka:9092"
	topic     = "celcius-readings"
	partition = 0
)

func main() {
	logger := log.NewLogfmtLogger(os.Stdout)
	logger = level.NewFilter(logger, level.AllowError())

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	level.Info(logger).Log("publisher setting up Kafka topic", topic)
	_, err := setupKafka(logger)
	if err != nil {
		level.Error(logger).Log("publisher failed to connect to Kafka topic", topic, "err", err)
		return
	}
	level.Info(logger).Log("publisher connected to Kafka topic", topic)

	go startPublishing(logger)
	<-sigchan
}

func setupKafka(logger log.Logger) (*kafka.Conn, error) {
	level.Info(logger).Log("connecting", address)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := kafka.DialLeader(ctx, "tcp", address, topic, partition)
	if err != nil {
		return nil, err
	}
	level.Info(logger).Log("connected", address)

	if err := conn.CreateTopics([]kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}...); err != nil {
		return nil, err
	}

	return conn, nil
}

func startPublishing(logger log.Logger) {
	rand.Seed(time.Now().Unix())
	w := &kafka.Writer{
		Addr:         kafka.TCP(address),
		Topic:        topic,
		RequiredAcks: kafka.RequireAll,
	}

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case t := <-ticker.C:
			publishTemperature(w, t, logger)
		}
	}
}

func publishTemperature(w *kafka.Writer, t time.Time, logger log.Logger) {
	id, err := uuid.NewUUID()
	if err != nil {
		level.Error(logger).Log("failed to create message ID", err)
	}
	msg := fmt.Sprintf(`{"temperature": %.2f, "timestamp":%d}`, rand.Float64()*100, t.Unix())
	rawMsg := json.RawMessage(msg)

	level.Info(logger).Log("publishing temperature reading", msg)
	err = w.WriteMessages(context.Background(), kafka.Message{
		Key:   id.NodeID(),
		Value: rawMsg,
	})
	if err != nil {
		level.Error(logger).Log("failed to publish message", id.String(), "error", err)
	}
}
