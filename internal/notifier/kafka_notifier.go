// internal/notifier/kafka_notifier.go
package notifier

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"

	"github.com/reactivejson/usr-svc/internal/domain"
)

// KafkaNotifier represents the Kafka notifier.
type KafkaNotifier struct {
	producer sarama.AsyncProducer
	topic    string
}

// NewKafkaNotifier creates a new instance of KafkaNotifier.
func NewKafkaNotifier(bootstrapServers, topic string) (*KafkaNotifier, error) {
	config := sarama.NewConfig()
	producer, err := sarama.NewAsyncProducer([]string{bootstrapServers}, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Kafka producer")
	}

	notifier := &KafkaNotifier{
		producer: producer,
		topic:    topic,
	}

	go notifier.handleSuccesses()
	go notifier.handleErrors()

	return notifier, nil
}

// NotifyUserChange sends a notification about a user change event.
func (n *KafkaNotifier) NotifyUserChange(user *domain.User) error {
	message := &sarama.ProducerMessage{
		Topic: n.topic,
		Value: sarama.StringEncoder(user.ID),
	}

	n.producer.Input() <- message

	return nil
}

// handleSuccesses handles successful message deliveries.
func (n *KafkaNotifier) handleSuccesses() {
	for range n.producer.Successes() {
		// Do nothing for now
	}
}

// handleErrors handles message delivery errors.
func (n *KafkaNotifier) handleErrors() {
	for err := range n.producer.Errors() {
		log.Printf("Kafka producer error: %v", err)
	}
}
