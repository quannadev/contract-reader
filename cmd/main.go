package main

import (
	"context"
	consumer2 "contract-reader/consumer"
	"github.com/KyberNetwork/logger"
)

func main() {
	//get config from args
	l, err := logger.InitLogger(logger.Configuration{
		EnableConsole:    true,
		ConsoleLevel:     "info",
		EnableJSONFormat: true,
	}, logger.LoggerBackend(1))
	if err != nil {
		return
	}
	ctx := context.Background()
	//create new consumer
	consumer := consumer2.NewNastConsumer(ctx, "topic")
	blockChan := make(chan consumer2.MessageBlock)
	//start consumer
	err = consumer.Start(ctx, blockChan)
	if err != nil {
		l.Fatalf("failed to start consumer: %v", err)
		return
	}
	//listen to blockChan
	for {
		select {
		case block := <-blockChan:
			l.Infof("received block: %v", block)
		}
	}
}
