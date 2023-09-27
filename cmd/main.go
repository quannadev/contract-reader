package main

import (
	"context"
	config2 "contract-reader/config"
	consumer2 "contract-reader/consumer"
	contract_reader "contract-reader/contract-reader"
	"contract-reader/producer"
	"contract-reader/utils"
	"github.com/KyberNetwork/logger"
	"github.com/ethereum/go-ethereum/common"
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
	config := config2.NewConfig() //get config from env

	//create new consumer
	consumer := consumer2.NewNastConsumer(ctx, config.Topic)
	//create new producer
	publisher := producer.NewNastProducer(config.Topic)

	//create new reader
	//todo init list contract from config or redis
	listContract := make(map[common.Address]utils.Contract)
	for _, contract := range config.GetListContracts() {
		listContract[contract.Address] = contract
	}
	reader := contract_reader.NewContractReader(publisher, listContract)
	//start consumer
	blockChan := make(chan utils.MessageBlock)
	go func() {
		err = consumer.Start(blockChan)
		if err != nil {
			l.Fatalf("failed to start consumer: %v", err)
			return
		}
	}()
	//start reader
	for {
		select {
		case block := <-blockChan:
			l.Infof("received block: %v", block)
			reader.OnMessage(ctx, &block)
		}
	}
}
