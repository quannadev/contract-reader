package contract_reader

import (
	"context"
	contractparser "contract-reader/contract-parser"
	"contract-reader/contract-parser/parsers"
	"contract-reader/producer"
	"contract-reader/protobuf/pb"
	"contract-reader/utils"
	"github.com/KyberNetwork/logger"
	"github.com/ethereum/go-ethereum/common"
	"time"
)

type ContractReader struct {
	producer    producer.IProducer
	listAddress map[common.Address]utils.ContractType
}

func (c *ContractReader) OnMessage(ctx context.Context, message *utils.MessageBlock) {
	ctxChild := context.WithoutCancel(ctx)
	logger.Infof("processing block: %d", message.Block.BlockNumber)
	block := message.Block
	events := pb.Events{
		Events:      make([]*pb.Event, 0),
		BlockNumber: block.BlockNumber,
		BlockHash:   common.Hex2Bytes(block.BlockHash),
		Timestamp:   uint64(time.Now().Unix()),
		ChainId:     block.ChainId,
	}
	for index, log := range block.Logs {
		logger.Infof("log %d: %v", index, log)
		contractType := c.getContractTypeFromAddress(common.HexToAddress(log.Address))
		if contractType != utils.UnSupportContract {
			parser := c.GetParser(contractType)
			event, err := parser.Parse(log)
			if err != nil {
				logger.Errorf("failed to parse event: %v", err)
				continue
			}
			events.Events = append(events.Events, event)
		}
	}
	logger.Infof("got events: %d from block %d", len(events.Events), message.Block.BlockNumber)
	err := c.producer.Publish(ctxChild, &events)
	if err != nil {
		//todo retry
		logger.Errorf("failed to publish message: %v", err)
	}
	logger.Infof("processed block: %d", message.Block.BlockNumber)
	message.Callback()
}

func (c *ContractReader) GetParser(contractType utils.ContractType) contractparser.IContractParser {
	switch contractType {
	case utils.ERC20:
		return parsers.NewErc20Parser()
	}
	logger.Fatalf("contract %s: un supported", contractType)
	panic("invalid contract type")
}

func (c *ContractReader) getContractTypeFromAddress(address common.Address) utils.ContractType {
	contactType, ok := c.listAddress[address]
	if !ok {
		return utils.UnSupportContract
	}
	return contactType
}

func NewContractReader(producer producer.IProducer, listAddress map[common.Address]utils.ContractType) IReader {
	return &ContractReader{
		producer:    producer,
		listAddress: listAddress,
	}
}
