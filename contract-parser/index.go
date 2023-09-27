package contract_parser

import (
	"contract-reader/protobuf/pb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func ToArrayHash(input []string) []common.Hash {
	arrHash := make([]common.Hash, len(input))
	for k, v := range input {
		arrHash[k] = common.HexToHash(v)
	}
	return arrHash
}

func ConvertLog(log *pb.Log) types.Log {
	return types.Log{
		Address:     common.HexToAddress(log.GetAddress()),
		Topics:      ToArrayHash(log.GetTopics()),
		Data:        common.Hex2Bytes(log.GetData()),
		BlockNumber: log.GetBlockNumber(),
		TxHash:      common.HexToHash(log.GetTransactionHash()),
		TxIndex:     uint(log.GetTransactionIndex()),
		BlockHash:   common.HexToHash(log.GetBlockHash()),
		Index:       uint(log.GetLogIndex()),
		Removed:     false,
	}
}
