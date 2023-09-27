package utils

import "contract-reader/protobuf/pb"

type MessageBlock struct {
	Block    *pb.Block
	Callback func()
}
