package parsers

import (
	contractparser "contract-reader/contract-parser"
	"contract-reader/protobuf/pb"
	"contract-reader/utils"
	"errors"
	"strings"
)

type Erc20Parser struct {
	contractparser.IParser
}

func (e *Erc20Parser) Parse(event *pb.Log) (*pb.Event, error) {
	eventName, err := e.GetEventName(event)
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(eventName) {
	case "transfer":
		return e.ParseEvent(event, new(Transfer))
	case "approval":
		return e.ParseEvent(event, new(Approval))
	}
	return nil, errors.New("event name unknown")
}

func NewErc20Parser() contractparser.IContractParser {
	return &Erc20Parser{
		contractparser.NewContractParser(utils.ERC20),
	}
}

type Approval struct {
	Owner   string `json:"owner"`
	Spender string `json:"spender"`
	Value   string `json:"value"`
}

type Transfer struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}
