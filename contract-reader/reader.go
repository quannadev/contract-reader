package contract_reader

import (
	"bytes"
	"context"
	"contract-reader/producer"
	"contract-reader/protobuf/pb"
	"contract-reader/utils"
	"github.com/KyberNetwork/logger"
	"github.com/ethereum/go-ethereum/common"
	"time"
)

type ContractReader struct {
	producer    producer.IProducer
	listAddress []common.Address
}

func (c *ContractReader) OnMessage(ctx context.Context, message *utils.MessageBlock) {
	//ctxChild := context.WithoutCancel(ctx)
	block := message.Block
	for index, log := range block.Logs {
		logger.Infof("log %d: %v", index, log)
		events := pb.Events{
			Events:      make([]*pb.Event, 0),
			BlockNumber: block.BlockNumber,
			BlockHash:   common.Hex2Bytes(block.BlockHash),
			Timestamp:   uint64(time.Now().Unix()),
			ChainId:     block.ChainId,
		}
		if c.checkAddress(common.HexToAddress(log.Address)) {
			//todo: parse log
		}
		err := c.producer.Publish(ctx, &events)
		if err != nil {
			//todo retry
			logger.Errorf("failed to publish message: %v", err)
		}
		message.Callback()
	}
}

func (c *ContractReader) checkAddress(address common.Address) bool {
	for _, addr := range c.listAddress {
		if bytes.Compare(address.Bytes(), addr.Bytes()) == 0 {
			return true
		}
	}
	return false
}
