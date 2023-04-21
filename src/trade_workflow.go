package app

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	protos "tradecapture-temporal-poc/app/src/generated"
)

const (
	activityName = "ProcessTradeActivity"
)

func TradeCaptureWorkflow(ctx workflow.Context, trade *protos.Trade, transport Transport) (*protos.TradeStatus, error) {

	// Configure retry policy for the activity and options
	activityOptions := workflow.ActivityOptions{
		TaskQueue:              TaskQueue,
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:        time.Second,
			BackoffCoefficient:     2.0,
			MaximumInterval:        100 * time.Second,
			MaximumAttempts:        0, // unlimited retries
			NonRetryableErrorTypes: []string{"InvalidTrade"},
		},
	}

	// apply the options
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// output of this trade workflow, need to be var because is modified by the activity result
	var tradeStatus protos.TradeStatus

	// execute the activity
	err := workflow.ExecuteActivity(ctx, ValidateTrade, transport, trade).Get(ctx, &tradeStatus)
	if err != nil {
		return nil, err
	}

	return &tradeStatus, nil
}
