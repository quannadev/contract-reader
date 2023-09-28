package contract_reader

import (
	"context"
	contractparser "contract-reader/contract-parser"
	"contract-reader/producer"
	"contract-reader/protobuf/pb"
	"contract-reader/utils"
	"errors"
	"github.com/KyberNetwork/logger"
	"github.com/ethereum/go-ethereum/common"
	"time"
)

type ContractReader struct {
	producer    producer.IProducer
	listAddress map[common.Address]utils.Contract
	listParser  map[common.Address]contractparser.IContractReader
}

func (c *ContractReader) OnMessage(ctx context.Context, message *utils.MessageBlock) {
	ctxChild := context.WithoutCancel(ctx)
	logger.Infof("processing block: %d", message.Block.BlockNumber)
	block := message.Block
	events := c.handleBlock(block)
	logger.Infof("got events: %d from block %d", len(events.Events), message.Block.BlockNumber)
	err := c.producer.Publish(ctxChild, events)
	if err != nil {
		//todo retry
		logger.Errorf("failed to publish message: %v", err)
	}
	logger.Infof("processed block: %d", message.Block.BlockNumber)
	message.Callback()
}

func (c *ContractReader) handleBlock(block *pb.Block) *pb.Events {
	events := pb.Events{
		Events:      make([]*pb.Event, 0),
		BlockNumber: block.BlockNumber,
		BlockHash:   common.Hex2Bytes(block.BlockHash),
		Timestamp:   uint64(time.Now().Unix()),
		ChainId:     block.ChainId,
	}
	for index, log := range block.Logs {
		logger.Infof("log %d: %v", index, log)
		logAddress := common.HexToAddress(log.Address)
		contract, err := c.getContractTypeFromAddress(logAddress)
		if err != nil {
			logger.Errorf("failed to get contract type from address: %v", err)
			continue
		}
		if contract.StartBlock > block.BlockNumber {
			logger.Debugf("contract: %s not start yet", contract.Address.String())
			continue
		}
		parser, ok := c.listParser[logAddress]
		if !ok {
			logger.Errorf("parser not found for contract: %s", contract.Address.String())
			continue
		}
		event, err := parser.Parse(log)
		if err != nil {
			logger.Errorf("failed to parse event: %v", err)
			continue
		}
		//todo handle event with new pool address => add to list address with this contract type
		events.Events = append(events.Events, event)
	}
	return &events
}

func (c *ContractReader) getContractTypeFromAddress(address common.Address) (*utils.Contract, error) {
	contact, ok := c.listAddress[address]
	if !ok {
		return nil, errors.New("contract not found")
	}
	return &contact, nil
}

func NewContractReader(producer producer.IProducer, listAddress map[common.Address]utils.Contract) IReader {
	listParser := make(map[common.Address]contractparser.IContractReader)
	for _, contract := range listAddress {
		listParser[contract.Address] = contractparser.NewEventReader(contract.AbiPath, contract.Events)
	}
	return &ContractReader{
		producer:    producer,
		listAddress: listAddress,
	}
}
