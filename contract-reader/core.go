package contract_reader

import (
	"context"
	"contract-reader/utils"
)

type IReader interface {
	OnMessage(ctx context.Context, message *utils.MessageBlock)
}
