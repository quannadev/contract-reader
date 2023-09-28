package contract_parser

import (
	abifiles "contract-reader/abi-files"
	"contract-reader/protobuf/pb"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KyberNetwork/logger"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/samber/lo"
)

type IContractReader interface {
	Parse(event *pb.Log) (*pb.Event, error)
	GetEventName(event *pb.Log) (*abi.Event, error)
	GetAbiContract() abi.ABI
}

type EventReader struct {
	abi.ABI
	events []string
}

func (c *EventReader) Parse(event *pb.Log) (*pb.Event, error) {
	eventLog, err := c.GetEventName(event)
	eventName := eventLog.Name
	if err != nil {
		return nil, err
	}
	_, exits := lo.Find(c.events, func(item string) bool {
		return item == eventName
	})
	if !exits {
		return nil, errors.New(fmt.Sprintf("event: %s not exits in contract", eventName))
	}
	logger.Debugf("parse event: %s", eventName)
	log := ConvertLog(event)
	result := make(map[string]interface{})
	indexed := make([]abi.Argument, 0)
	for _, arg := range eventLog.Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	//if err := abi.ParseTopicsIntoMap(result, indexed, log.Topics[1:]); err != nil {
	//	logger.Errorf("failed to parse topics into map: %v", err)
	//	return nil, err
	//}
	fmt.Printf("result: %v\n", result)
	if err := c.UnpackIntoMap(result, eventName, log.Data); err != nil {
		logger.Errorf("failed to unpack data into map: %v", err)
		return nil, err
	}

	packedData, err := json.Marshal(result)
	if err != nil {
		logger.Fatalf("failed to marshal data: %v", err)
		//panic application when marshal data failed
		return nil, err
	}
	return &pb.Event{
		LogIndex:  event.GetLogIndex(),
		TxHash:    event.GetTransactionHash(),
		Address:   event.GetAddress(),
		Data:      packedData,
		Extra:     nil, //todo call rpc get extra data
		EventName: eventName,
	}, nil
}

func NewEventReader(abiPath string, events []string) IContractReader {
	contract := abifiles.GetAbiContract(abiPath)
	return &EventReader{
		contract,
		events,
	}
}

func (c *EventReader) GetAbiContract() abi.ABI {
	return c.ABI
}

func (c *EventReader) GetEventName(event *pb.Log) (*abi.Event, error) {
	eventLog, err := c.EventByID(common.HexToHash(event.GetTopics()[0]))
	if err != nil {
		return nil, err
	}
	return eventLog, nil
}
