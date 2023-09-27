package contract_parser

import (
	abifiles "contract-reader/abi-files"
	"contract-reader/protobuf/pb"
	"contract-reader/utils"
	"encoding/json"
	"github.com/KyberNetwork/logger"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type IContractParser interface {
	Parse(event *pb.Log) (*pb.Event, error)
}

type IParser interface {
	ParseEvent(event *pb.Log, output interface{}) (*pb.Event, error)
	GetEventName(event *pb.Log) (string, error)
	GetAbiContract() abi.ABI
}

type ContractParser struct {
	contract abi.ABI
}

func (c *ContractParser) ParseEvent(event *pb.Log, data interface{}) (*pb.Event, error) {
	eventLog, err := c.contract.EventByID(common.HexToHash(event.GetTopics()[0]))
	if err != nil {
		return nil, err
	}
	logger.Infof("parse event: %s", eventLog.Name)
	contract := bind.NewBoundContract(common.HexToAddress(event.Address), c.contract, nil, nil, nil)
	log := ConvertLog(event)
	if err := contract.UnpackLog(data, eventLog.Name, log); err != nil {
		logger.Errorf("unpack log error: %v", err)
		return nil, err
	}
	packedData, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("failed to marshal data: %v", err)
		return nil, err
	}
	return &pb.Event{
		LogIndex:  event.GetLogIndex(),
		TxHash:    event.GetTransactionHash(),
		Address:   event.GetAddress(),
		Data:      packedData,
		Extra:     nil,
		EventName: eventLog.Name,
	}, nil
}

func NewContractParser(contractType utils.ContractType) IParser {
	contract := abifiles.GetAbiContract(contractType)
	return &ContractParser{
		contract: contract,
	}
}

func (c *ContractParser) GetAbiContract() abi.ABI {
	return c.contract
}

func (c *ContractParser) GetEventName(event *pb.Log) (string, error) {
	eventLog, err := c.contract.EventByID(common.HexToHash(event.GetTopics()[0]))
	if err != nil {
		return "", err
	}
	return eventLog.Name, nil
}
