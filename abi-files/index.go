package abi_files

import (
	"bytes"
	"contract-reader/utils"
	_ "embed"
	"github.com/KyberNetwork/logger"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

func GetAbiContract(abiPath string) abi.ABI {
	value, err := utils.ReadFile(abiPath)
	var abiContract abi.ABI
	abiContract, err = abi.JSON(bytes.NewReader(value))
	if err != nil {
		logger.Fatalf("failed to get abi for contract %s: %v", abiPath, err)
	}
	return abiContract
}
