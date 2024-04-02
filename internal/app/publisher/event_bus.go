package publisher

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto" //nolint:staticcheck
	eventBus "github.com/voonik/goConnect/api/go/event_bus/publisher"
	eventBusPublisher "github.com/voonik/goConnect/event_bus/publisher"
	"google.golang.org/grpc"
)

var protoMarshal = proto.Marshal

var EventBusClient = func() eventBus.PublisherClient {
	return eventBusPublisher.Publisher()
}

type PublisherClient interface {
	Publish(ctx context.Context, in *eventBus.PublishRequest, opts ...grpc.CallOption) (*eventBus.PublishResponse, error)
}

func Publish(ctx context.Context, topic string, key proto.Message, value proto.Message) (*eventBus.PublishResponse, error) {
	marshalledKey, err := protoMarshal(key)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal key with error: %s", err.Error())
	}

	marshalledValue, err := protoMarshal(value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal value with error: %s", err.Error())
	}

	request := &eventBus.PublishRequest{
		Topic: topic,
		Key:   marshalledKey,
		Value: marshalledValue,
	}

	response, err := EventBusClient().Publish(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("[EventBus] publish failed with error: %v", err)
	}

	return response, nil
}
