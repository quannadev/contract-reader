package consumer

import (
	"context"
	"contract-reader/protobuf/pb"
	"contract-reader/utils"
	"encoding/json"
	"github.com/KyberNetwork/logger"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"log"
	"time"
)

type NastConsumer struct {
	consumer jetstream.Consumer
	topic    string
}

func (n *NastConsumer) Start(onMessage chan utils.MessageBlock) error {
	logger.Infof("start consumer with topic: %v", n.topic)
	messagesContext, err := n.consumer.Messages()
	if err != nil {
		logger.Errorf("failed to get messages: %v", err)
	}
	defer messagesContext.Stop()
	for {
		select {
		default:
			msg, err := messagesContext.Next()
			if err != nil {
				logger.Errorf("failed to get next message: %v", err)
			}
			logger.Infof("received message: %v", string(msg.Data()))
			block := pb.Block{}
			err = json.Unmarshal(msg.Data(), &block)
			if err != nil {
				logger.Errorf("failed to unmarshal message: %v", err)
				return err
			}
			onMessage <- utils.MessageBlock{
				Block: &block,
				Callback: func() {
					logger.Infof("ack message: %v", string(msg.Data()))
					err := msg.Ack()
					if err != nil {
						logger.Errorf("failed to ack message: %v", err)
					}
				},
			}
		}
	}
}

func NewNastConsumer(ctx context.Context, topic string) IConsumer {
	client, err := nats.Connect(nats.DefaultURL, setupConnOptions([]nats.Option{})...)
	if err != nil {
		logger.Fatalf("failed to connect to nats server: %v", err)
	}
	logger.Infof("connected to nats server: %v", nats.DefaultURL)
	ctxChild, cancel := context.WithTimeout(ctx, 60*time.Minute)
	defer cancel()
	stream, err := jetstream.New(client)
	if err != nil {
		logger.Errorf("failed to create stream: %v", err)
	}
	s, err := stream.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "block_consumer",
		Subjects: []string{topic},
	})
	if err != nil {
		logger.Fatalf("failed to create stream: %v", err)
	}
	consumer, err := s.OrderedConsumer(ctxChild, jetstream.OrderedConsumerConfig{
		MaxResetAttempts: 5,
	})
	if err != nil {
		logger.Errorf("failed to create consumer: %v", err)
	}
	return &NastConsumer{
		consumer: consumer,
	}
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second
	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}
