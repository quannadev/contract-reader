package consumer

import (
	"context"
	"contract-reader/producer"
	"fmt"
	"github.com/KyberNetwork/logger"
	"testing"
	"time"
)

func TestNastConsumer_Start(t *testing.T) {
	_, err := logger.InitLogger(logger.Configuration{
		EnableConsole:    true,
		ConsoleLevel:     "debug",
		EnableJSONFormat: true,
	}, logger.LoggerBackend(1))
	if err != nil {
		return
	}
	topic := "topic"
	ctx := context.Background()
	consumer := NewNastConsumer(ctx, topic)
	publisher := producer.NewNastProducer(topic)
	//create a publisher to send messages to a topic with go routine
	go func() {
		//start consumer
		t.Logf("start creating consumer")
		msgChan := make(chan MessageBlock)
		go func() {
			err = consumer.Start(msgChan)
			if err != nil {
				t.Errorf("failed to start consumer: %v", err)
				return
			}
		}()
		t.Logf("start consumer successfully")
		//listen to msgChan
		counter := 0
		for {
			select {
			case msg := <-msgChan:
				counter++
				//ack message
				msg.callback()
				//if counter == 100 {
				//	return
				//}
			}
		}
	}()

	//send 100 messages to a topic
	for i := 0; i < 100; i++ {
		time.Sleep(1 * time.Second)
		err := publisher.Publish(ctx, fmt.Sprintf("message %v", i))
		if err != nil {
			t.Errorf("failed to publish message: %v", err)
		}
	}
}
