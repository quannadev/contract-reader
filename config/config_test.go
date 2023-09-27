package config

import (
	"contract-reader/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	//set env
	os.Setenv("CONTRACT_PATH", "contract.json")
	os.Setenv("TOPIC", "test")
	os.Setenv("CHAIN_ID", "1")
	config := NewConfig()
	assert.Equal(t, config.ContractPath, "contract.json")
	assert.Equal(t, config.Topic, "test")
	assert.Equal(t, config.ChainId, int64(1))
	contract := config.GetListContracts()
	assert.Equal(t, len(contract), 1)
	assert.Equal(t, contract[0].Address.String(), "0xdac17f958d2ee523a2206206994597c13d831ec7")
	assert.Equal(t, contract[0].Type, utils.ERC20)
	assert.Equal(t, contract[0].StartBlock, uint64(0))
}
