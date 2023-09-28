package utils

import (
	"contract-reader/protobuf/pb"
	"github.com/ethereum/go-ethereum/common"
)

type MessageBlock struct {
	Block    *pb.Block
	Callback func()
}

type Contract struct {
	Address    common.Address `json:"address"`
	Type       string         `json:"type"`
	StartBlock uint64         `json:"start_block"`
	AbiPath    string         `json:"abi_path"`
	Events     []string       `json:"events"`
}
