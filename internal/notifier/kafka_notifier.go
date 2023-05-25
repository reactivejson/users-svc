package notifier

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"log"

	"github.com/reactivejson/users-svc/internal/domain"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2023
 */

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
/*func (n *KafkaNotifier) NotifyUserChange(user *domain.User) error {
	message := &sarama.ProducerMessage{
		Topic: n.topic,
		Value: sarama.StringEncoder(user.ID),
	}

	n.producer.Input() <- message

	return nil
}*/

func (n *KafkaNotifier) NotifyUserChange(eventType UserEventType, user *domain.User) error {
	payload, err := json.Marshal(UserEvent{UserID: user.ID, Type: eventType, User: user})
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: n.topic,
		Value: sarama.StringEncoder(payload),
	}
	n.producer.Input() <- msg

	return nil
}

func (n *KafkaNotifier) Close() error {
	return n.producer.Close()
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

type UserEventType string

var (
	Created UserEventType = "CREATED"
	Updated UserEventType = "UPDATED"
	Deleted UserEventType = "DELETED"
)

type UserEvent struct {
	Type   UserEventType `json:"type"`
	UserID string        `json:"user_id"`
	User   *domain.User  `json:"user,omitempty"`
}
