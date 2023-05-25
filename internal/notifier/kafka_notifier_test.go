//go:build integrations
// +build integrations

package notifier

import (
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/reactivejson/users-svc/internal/domain"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2023
 */

type MockAsyncProducer struct {
	mock.Mock
}

func (m *MockAsyncProducer) AsyncClose() {
	//TODO implement me
	panic("implement me")
}

func (m *MockAsyncProducer) Close() error {
	//TODO implement me
	panic("implement me")
}

func (m *MockAsyncProducer) IsTransactional() bool {
	//TODO implement me
	panic("implement me")
}

func (m *MockAsyncProducer) TxnStatus() sarama.ProducerTxnStatusFlag {
	//TODO implement me
	panic("implement me")
}

func (m *MockAsyncProducer) BeginTxn() error {
	//TODO implement me
	panic("implement me")
}

func (m *MockAsyncProducer) CommitTxn() error {
	//TODO implement me
	panic("implement me")
}

func (m *MockAsyncProducer) AbortTxn() error {
	//TODO implement me
	panic("implement me")
}

func (m *MockAsyncProducer) AddOffsetsToTxn(offsets map[string][]*sarama.PartitionOffsetMetadata, groupId string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockAsyncProducer) AddMessageToTxn(msg *sarama.ConsumerMessage, groupId string, metadata *string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockAsyncProducer) Input() chan<- *sarama.ProducerMessage {
	args := m.Called()
	return args.Get(0).(chan *sarama.ProducerMessage)
}

func (m *MockAsyncProducer) Successes() <-chan *sarama.ProducerMessage {
	args := m.Called()
	return args.Get(0).(<-chan *sarama.ProducerMessage)
}

func (m *MockAsyncProducer) Errors() <-chan *sarama.ProducerError {
	args := m.Called()
	return args.Get(0).(<-chan *sarama.ProducerError)
}

func TestNotifyUserChange(t *testing.T) {
	// Create a mock AsyncProducer
	mockProducer := new(MockAsyncProducer)

	// Set up expectations for Input method
	inputChan := make(chan *sarama.ProducerMessage)
	mockProducer.On("Input").Return(inputChan)

	// Create KafkaNotifier with the mock producer
	notifier := &KafkaNotifier{
		producer: mockProducer,
		topic:    "user-events",
	}

	// Create a user to notify
	user := &domain.User{
		ID:        "1",
		FirstName: "Alice",
		LastName:  "Bob",
		Country:   "UK",
	}

	// Trigger the notification
	err := notifier.NotifyUserChange(Created, user)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Expect the message to be enqueued to the producer's input channel
	expectedMessage := &sarama.ProducerMessage{
		Topic: "user-events",
		Value: sarama.StringEncoder("1"),
	}
	select {
	case msg := <-inputChan:
		assert.Equal(t, expectedMessage, msg)
	case <-time.After(time.Second):
		t.Error("Expected message not sent to producer")
	}

	// Assert that all expectations were met
	mockProducer.AssertExpectations(t)
}

func TestHandleSuccesses(t *testing.T) {
	// Create a mock AsyncProducer
	mockProducer := new(MockAsyncProducer)

	// Set up expectations for Successes method
	successesChan := make(chan *sarama.ProducerMessage)
	mockProducer.On("Successes").Return(successesChan)

	// Create KafkaNotifier with the mock producer
	notifier := &KafkaNotifier{
		producer: mockProducer,
		topic:    "user-events",
	}

	// Start the handleSuccesses goroutine
	go notifier.handleSuccesses()

	// Simulate a successful message delivery
	successMessage := &sarama.ProducerMessage{Topic: "user-events", Value: sarama.StringEncoder("1")}
	successesChan <- successMessage

	// Wait for the handleSuccesses goroutine to handle the message
	time.Sleep(100 * time.Millisecond)

	// Assert that all expectations were met
	mockProducer.AssertExpectations(t)
}

func TestHandleErrors(t *testing.T) {
	// Create a mock AsyncProducer
	mockProducer := new(MockAsyncProducer)

	// Set up expectations for Errors method
	errorsChan := make(chan *sarama.ProducerError)
	mockProducer.On("Errors").Return(errorsChan)

	// Create KafkaNotifier with the mock producer
	notifier := &KafkaNotifier{
		producer: mockProducer,
		topic:    "user-events",
	}

	// Start the handleErrors goroutine
	go notifier.handleErrors()

	// Simulate a message delivery error
	errorMessage := &sarama.ProducerError{
		Err: sarama.ErrOutOfBrokers,
		Msg: &sarama.ProducerMessage{Topic: "user-events", Value: sarama.StringEncoder("1")},
	}
	errorsChan <- errorMessage

	// Wait for the handleErrors goroutine to handle the error
	time.Sleep(100 * time.Millisecond)

	// Assert that all expectations were met
	mockProducer.AssertExpectations(t)
}
