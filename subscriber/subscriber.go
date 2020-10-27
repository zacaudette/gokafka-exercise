package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
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

	level.Info(logger).Log("subscriber listening for readings", topic)
	go readings(logger)
	<-sigchan
}

func readings(logger log.Logger) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{address},
		Topic:     topic,
		Partition: partition,
	})

	for {
		message, err := r.ReadMessage(context.Background())
		if err != nil {
			level.Error(logger).Log("failed to read message", err)
			continue
		}
		level.Info(logger).Log("message received", string(message.Key), "time", message.Time.String())
		fmt.Println(string(message.Value))
	}

}
