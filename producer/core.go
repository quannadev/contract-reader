package producer

import (
	"context"
	"contract-reader/protobuf/pb"
)

type IProducer interface {
	Publish(ctx context.Context, events *pb.Events) error
}
