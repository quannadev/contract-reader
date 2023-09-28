package contract_parser

import (
	"contract-reader/protobuf/pb"
	"testing"
)

// url: https://etherscan.io/tx/0x04696778c2a3f7154a9a61f1145be0e8e81ac24db61e1b96ec9b75bc172330d4#eventlog
func TestERC20(t *testing.T) {
	blockNumber := uint64(18229689)
	txHash := "0x04696778c2a3f7154a9a61f1145be0e8e81ac24db61e1b96ec9b75bc172330d4"
	logIndex := uint64(304)
	log := pb.Log{
		Address: "0xdac17f958d2ee523a2206206994597c13d831ec7",
		Topics: []string{
			"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
			"0x0000000000000000000000002008b6c3d07b061a84f790c035c2f6dc11a0be70",
			"0x00000000000000000000000022f9dcf4647084d6c31b2765f6910cd85c178c18",
		},
		Data:            "0x000000000000000000000000000000000000000000000000000000000e3c6f7b",
		BlockNumber:     &blockNumber,
		TransactionHash: &txHash,
		LogIndex:        &logIndex,
	}
	contract := NewEventReader("abi-files/ERC20.json", []string{"Transfer"})
	event, err := contract.Parse(&log)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(event)
}
