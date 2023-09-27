package abi_files

import (
	"bytes"
	"contract-reader/utils"
	_ "embed"
	"github.com/KyberNetwork/logger"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

//go:embed ERC20.json
var ERC20 []byte

func GetAbiContract(contractType utils.ContractType) abi.ABI {
	var abiContent []byte
	switch contractType {
	case utils.ERC20:
		abiContent = ERC20
		break
	default:
		logger.Fatalf("contract %s: un supported", contractType)
	}
	var abiContract abi.ABI
	abiContract, err := abi.JSON(bytes.NewReader(abiContent))
	if err != nil {
		logger.Fatalf("failed to get abi for contract %s: %v", contractType, err)
	}
	return abiContract
}
