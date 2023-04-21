package app

import (
	"context"

	"google.golang.org/protobuf/proto"
)

/*****
The idea of this is to have a transport layer interface that can have many switcheable implementations, kafka, grpc, rest, etc.
*****/

// payload input used by the transport layers
type MessagePayload struct {
	Key   string
	Value []byte
}

// interface for the transport be able to receive a message
type MessageReceiver interface {
	Receive(ctx context.Context) (*MessagePayload, error)
}

// interface for the transport be able to send a message
type MessageSender interface {
	Send(ctx context.Context, payload *MessagePayload) error
}

// Interface to be able to implement different transport layers, kafka, grpc, rest, etc.
type Transport interface {
	MessageReceiver
	MessageSender
}

// serializer and deserializer
func EncodeProtoMessage(message proto.Message) ([]byte, error) {
	return proto.Marshal(message)
}

func DecodeProtoMessage(data []byte, message proto.Message) error {
	return proto.Unmarshal(data, message)
}
