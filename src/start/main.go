package main

import (
	"context"
	"fmt"
	"log"
	"time"

	app "tradecapture-temporal-poc/app/src"
	protos "tradecapture-temporal-poc/app/src/generated"

	"github.com/segmentio/ksuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	ctx := context.Background()

	// Start Temporal worker
	workerOptions := worker.Options{
		MaxConcurrentActivityExecutionSize: app.MaxConns,
		BackgroundActivityContext:          context.Background(),
	}

	// Configure Temporal client
	//client, err := client.NewClient(client.Options{})
	//if err != nil {
	//	log.Fatalln("Unable to create client", err)
	//}

	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	// create a new worker
	worker := worker.New(c, app.TaskQueue, workerOptions)
	defer worker.Stop()

	// Register TradeCaptureWorkflow with worker
	worker.RegisterWorkflow(app.TradeCaptureWorkflow)
	//worker.RegisterActivity(ProcessTradeActivity)

	// Start Temporal worker
	err = worker.Start()
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}

	// Initialize Kafka transport
	kafkaTransport := app.NewKafkaTransport(app.IncomingTopic, app.OutgoingTopic, app.KafkaAddress)

	// Start message processor with workflow
	go processKafkaMessages(ctx, kafkaTransport, c)

	// Send test trades to incoming topic
	//sendTestTrades(kafkaTransport, 10)

	// Wait for user input to exit
	fmt.Println("Press Enter to exit...")
	var input string
	fmt.Scanln(&input)
}

// for test proposes
func sendTestTrades(transport app.Transport, count int) {
	for i := 0; i < count; i++ {
		trade := protos.Trade{
			Id:       ksuid.New().String(),
			Amount:   1000,
			Quantity: 5,
			Buyer:    "buyer1",
			Seller:   "seller1",
			Total:    5000,
		}
		data, err := app.EncodeProtoMessage(&trade)
		if err != nil {
			log.Printf("Error encoding trade: %v", err)
			continue
		}

		payload := &app.MessagePayload{
			Key:   trade.Id,
			Value: data,
		}

		err = transport.Send(context.Background(), payload)
		if err != nil {
			log.Printf("Error sending trade: %v", err)
			continue
		}

		fmt.Printf("Sent trade with ID %s\n", trade.Id)
	}
}

func processKafkaMessages(ctx context.Context, transport app.Transport, c client.Client) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			payload, err := transport.Receive(ctx)
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				continue
			}

			var trade protos.Trade
			err = app.DecodeProtoMessage(payload.Value, &trade)
			if err != nil {
				log.Printf("Error decoding message: %v", err)
				continue
			}

			// Call the TradeCaptureWorkflow
			workflowOptions := client.StartWorkflowOptions{
				ID:                       trade.Id,
				TaskQueue:                app.TaskQueue,
				WorkflowExecutionTimeout: 10 * time.Minute,
				WorkflowTaskTimeout:      time.Minute,
				//WorkflowIDReusePolicy:          client.WorkflowIDReusePolicyAllowDuplicate,
				//DisableWorkflowExecutionCancel: true,
			}
			we, err := c.ExecuteWorkflow(ctx, workflowOptions, app.TradeCaptureWorkflow, trade, transport)
			if err != nil {
				log.Printf("Error starting workflow: %v", err)
			} else {
				log.Printf("Started workflow for trade with ID: %s, runID: %s", trade.Id, we.GetRunID())
			}
		}
	}
}
