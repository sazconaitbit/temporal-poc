package app

import (
	"context"
	"errors"
	protos "tradecapture-temporal-poc/app/src/generated"

	"github.com/golang/protobuf/proto"
)

func ValidateTrade(ctx context.Context, transport Transport, input *protos.Trade) (*protos.TradeStatus, error) {
	// Validaciones
	if input.Quantity <= 0 || input.Amount <= 0 {
		return nil, errors.New("quantity and amount must be greater than 0")
	}
	if input.Total != input.Quantity*input.Amount {
		return nil, errors.New("total must be equal to quantity x amount")
	}

	// Producir el mensaje TradeStatus en el transporte
	tradeStatus := &protos.TradeStatus{
		ReferenceId: input.Id,
		Status:      protos.TradeStatus_ACCEPTED,
	}

	tradeStatusBytes, err := proto.Marshal(tradeStatus)
	if err != nil {
		return nil, err
	}

	err = transport.Send(ctx, &MessagePayload{
		Key:   input.Id,
		Value: tradeStatusBytes,
	})

	if err != nil {
		return nil, err
	}

	return tradeStatus, nil
}
