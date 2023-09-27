package consumer

import "contract-reader/utils"

type IConsumer interface {
	Start(onMessage chan utils.MessageBlock) error
}
