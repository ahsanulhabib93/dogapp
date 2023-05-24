package publisher

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/shopuptech/event-bus-logs-go/core"
	"github.com/shopuptech/event-bus-logs-go/ss2"
	"github.com/stretchr/testify/assert"
	eventBus "github.com/voonik/goConnect/api/go/event_bus/publisher"
	"github.com/voonik/ss2/internal/app/publisher/mocks"
)

func TestPublish(t *testing.T) {
	key, value := getLog(t)
	marshalledKey, _ := proto.Marshal(key)
	marshalledValue, _ := proto.Marshal(value)
	topic := "test-topic"
	ctx := context.Background()

	t.Run("should not panic when ctx, key and value are provided and publish to kafka", func(t *testing.T) {
		mockedEventBus, resetEventBus := mocks.SetupMockPublisherClient(t, &eventBusClient)
		defer resetEventBus()

		request := &eventBus.PublishRequest{
			Topic: topic,
			Key:   marshalledKey,
			Value: marshalledValue,
		}

		mockedEventBus.On("Publish", ctx, request).Return(&eventBus.PublishResponse{Success: true}, nil)

		response, err := Publish(ctx, topic, key, value)

		assert.NotNil(t, response)
		assert.True(t, response.Success)
		assert.NoError(t, err)
		mockedEventBus.AssertExpectations(t)
	})

	t.Run("should not publish to kafka in case of error", func(t *testing.T) {
		t.Run("when event-bus returns error", func(t *testing.T) {
			mockedEventBus, resetEventBus := mocks.SetupMockPublisherClient(t, &eventBusClient)
			defer resetEventBus()

			request := &eventBus.PublishRequest{
				Topic: topic,
				Key:   marshalledKey,
				Value: marshalledValue,
			}

			mockedEventBus.On("Publish", ctx, request).Return(nil, errors.New("cannot publish to kafka"))

			res, err := Publish(ctx, topic, key, value)

			assert.Nil(t, res)
			assert.EqualError(t, err, "[EventBus] publish failed with error: cannot publish to kafka")
			mockedEventBus.AssertExpectations(t)
		})

		t.Run("failed to marshal key", func(t *testing.T) {
			protoMarshal = func(m proto.Message) ([]byte, error) {
				return nil, fmt.Errorf("proto marshal key error")
			}
			defer func() { protoMarshal = proto.Marshal }()

			res, err := Publish(ctx, topic, key, value)

			assert.Nil(t, res)
			assert.EqualError(t, err, "failed to marshal key with error: proto marshal key error")
		})
	})
}

func getLog(t *testing.T) (*ss2.SupplierLogKey, *ss2.SupplierLogValue) {
	key := &ss2.SupplierLogKey{
		Event: &core.Event{
			Id: "test-id",
		},
		SupplierId: 123,
	}

	value := &ss2.SupplierLogValue{
		Id:     123,
		Name:   "supplier",
		Status: "created",
	}

	return key, value
}
