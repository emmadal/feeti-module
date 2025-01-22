package stream

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// KafkaProducer produce messages to kafka topic
func KafkaProducer(msg interface{}, topic, servers string) {

	// Error handling
	defer func() {
		if r := recover(); r != nil {
			logrus.WithFields(logrus.Fields{"error": r}).Error("Fatal error on kafka producer")
		}
	}()

	// Create producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": servers,
		"client.id":         "feeti-client",
	})

	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					logrus.WithFields(logrus.Fields{"error": ev.TopicPartition.Error}).Error(ev.TopicPartition)
				}
			}
		}
	}()

	// marshal message to bytes
	jsonMsg, _ := json.Marshal(msg)

	// Produce message
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          jsonMsg,
		Key:            []byte(uuid.NewString()),
	}, nil)

	// Flush messages to kafka
	p.Flush(1500)

	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("Error Message to send on kafka")
	} else {
		logrus.Info("Message sent to kafka successfully")
	}
}
