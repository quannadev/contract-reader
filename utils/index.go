package utils

import (
	"contract-reader/protobuf/pb"
)

type MessageBlock struct {
	Block    *pb.Block
	Callback func()
}

type ContractType string

const (
	ERC20             ContractType = "ERC20"
	UnSupportContract ContractType = "UnSupportContract"
)
