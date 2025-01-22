package stream

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

func KafkaConsumer[T any](servers string) {
	// Error handling
	defer func() {
		if r := recover(); r != nil {
			logrus.WithFields(logrus.Fields{"error": r}).Error("Fatal error on kafka consumer")
		}
	}()

	// Create consumer
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": servers,
		"group.id":          "feeti-group",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	// Topics subscription
	err = c.SubscribeTopics([]string{"payment", "campaign", "refund", "transfer", "system"}, nil)
	if err != nil {
		panic(err)
	}

	defer c.Close()

	// Read messages
	for {
		msg, err := c.ReadMessage(time.Millisecond)
		if err == nil {
			switch *msg.TopicPartition.Topic {
			case "payment":
				go payment[T](msg.Value)
				break
			case "campaign":
				go campaign[T](msg.Value)
				break
			case "refund":
				go refund[T](msg.Value)
				break
			case "transfer":
				go transfer[T](msg.Value)
				break
			case "system":
				go system[T](msg.Value)
				break
			default:
				break
			}
		}
	}
}
