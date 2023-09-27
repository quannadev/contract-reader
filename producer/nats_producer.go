package producer

import (
	"context"
	"contract-reader/protobuf/pb"
	"encoding/json"
	"github.com/KyberNetwork/logger"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NastProducer struct {
	stream jetstream.JetStream
	topic  string
}

func (n *NastProducer) Publish(ctx context.Context, events *pb.Events) error {
	message, err := json.Marshal(events)
	if err != nil {
		logger.Errorf("failed to marshal block: %v", err)
		return err
	}
	_, err = n.stream.Publish(ctx, n.topic, message)
	if err != nil {
		return err
	}
	logger.Infof("published message: %v", string(message))
	return nil
}

func NewNastProducer(topic string) IProducer {
	client, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		logger.Fatalf("failed to connect to nats server: %v", err)
	}
	logger.Infof("connected to nats server: %s", nats.DefaultURL)
	stream, err := jetstream.New(client)
	if err != nil {
		logger.Fatalf("failed to create stream: %v", err)
	}
	return &NastProducer{
		stream: stream,
		topic:  topic,
	}
}
