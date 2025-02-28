package bwork

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/tuanuet/bwork/proto" // Adjust import path
	"google.golang.org/grpc"
)

type Record struct {
	Data string
}

type Event struct {
	Data     string
	IsActive bool
}

type Client struct {
	nodeID         string
	client         pb.ShardServiceClient
	currentRecords []string
}

func NewNodeClient(target, nodeID string) (*Client, error) {
	conn, err := grpc.NewClient(target, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to dial ShardManager: %v", err)
	}
	return &Client{
		nodeID: nodeID,
		client: pb.NewShardServiceClient(conn),
	}, nil
}

func (nc *Client) Register(ctx context.Context) error {
	resp, err := nc.client.RegisterNode(ctx, &pb.NodeRequest{NodeId: nc.nodeID})
	if err != nil {
		return fmt.Errorf("node %s failed to register: %w", nc.nodeID, err)
	}
	if !resp.Success {
		return fmt.Errorf("node %s registration failed: %s", nc.nodeID, resp.Message)
	}
	log.Printf("Node %s registered successfully\n", nc.nodeID)
	return nil
}

func (nc *Client) FetchRecords(ctx context.Context) ([]*Record, error) {
	res, err := nc.client.GetRecordsForNode(ctx, &pb.NodeRequest{NodeId: nc.nodeID})
	if err != nil {
		log.Fatalf("Node %s failed to subscribe: %v", nc.nodeID, err)
	}

	records := make([]*Record, len(res.GetRecords()))
	for i, record := range res.Records {
		records[i] = &Record{
			Data: record.RecordId,
		}
	}
	return records, nil
}

func (nc *Client) Subscribe(ctx context.Context) chan []*Record {
	recordChan := make(chan []*Record)

	doFunc := func() error {
		stream, err := nc.client.SubscribeRebalance(ctx, &pb.NodeRequest{NodeId: nc.nodeID})
		if err != nil {
			return fmt.Errorf("node %s failed to subscribe: %w", nc.nodeID, err)
		}

		for {
			_, err = stream.Recv()
			if err != nil {
				return fmt.Errorf("node %s stream closed: %v", nc.nodeID, err)
			}

			records, err := nc.FetchRecords(ctx)
			if err != nil {
				return fmt.Errorf("node %s failed to subscribe: %v", nc.nodeID, err)
			} else {
				recordChan <- records
			}
		}
	}

	go func() {
		defer close(recordChan)
		for {
			if err := doFunc(); err != nil {
				log.Printf("Node %s failed to subscribe: %v", nc.nodeID, err)
				time.Sleep(1 * time.Second)
			}
		}
	}()

	return recordChan

}

func (nc *Client) SubscribeChanges(ctx context.Context) chan []*Event {
	recordChan := make(chan []*Event)

	doFunc := func() error {
		var stream pb.ShardService_SubscribeRebalanceClient
		var err error
		stream, err = nc.client.SubscribeRebalance(ctx, &pb.NodeRequest{NodeId: nc.nodeID})
		if err != nil {
			return fmt.Errorf("node %s failed to subscribe: %w", nc.nodeID, err)
		}

		log.Printf("Node %s subscribed to changes\n", nc.nodeID)

		for {
			_, err = stream.Recv()
			if err != nil {
				return fmt.Errorf("node %s stream closed: %v", nc.nodeID, err)
			}

			res, err := nc.client.GetRecordsForNode(ctx, &pb.NodeRequest{NodeId: nc.nodeID})
			if err != nil {
				log.Fatalf("Node %s failed to subscribe: %v", nc.nodeID, err)
			}

			recordIds := make([]string, len(res.GetRecords()))
			records := make([]*Event, len(res.GetRecords()))
			for i, record := range res.Records {
				recordIds[i] = record.GetRecordId()
				records[i] = &Event{
					Data:     record.GetRecordId(),
					IsActive: true,
				}
			}

			// delete = old - new
			deletes := DifferentArrays(nc.currentRecords, records, oldFunc, newFunc)
			// append = new - delete
			news := DifferentArrays(records, nc.currentRecords, newFunc, oldFunc)

			events := make([]*Event, 0)
			for _, n := range news {
				events = append(events, &Event{
					Data:     n.Data,
					IsActive: true,
				})
			}
			for _, d := range deletes {
				events = append(events, &Event{
					Data:     d,
					IsActive: false,
				})
			}

			recordChan <- events
			nc.currentRecords = recordIds
		}

		return nil
	}

	go func() {
		defer close(recordChan)
		for {
			if err := doFunc(); err != nil {
				log.Printf("Node %s failed to subscribe: %v", nc.nodeID, err)
				time.Sleep(1 * time.Second)
			}
		}
	}()
	return recordChan

}

func (nc *Client) StartPinging(ctx context.Context) {
	errChan := make(chan error)
	go func() {
		defer close(errChan)
		for {
			ctxTimeout, cancel := context.WithTimeout(ctx, time.Second)
			resp, err := nc.client.Ping(ctxTimeout, &pb.PingRequest{NodeId: nc.nodeID})
			if err != nil || !resp.Success {
				errChan <- fmt.Errorf("node %s ping failed: %w", nc.nodeID, err)
			}
			cancel()
			time.Sleep(2 * time.Second) // Ping every 3 seconds
		}
	}()

	go func() {
		for err := range errChan {
			if err = nc.Register(ctx); err != nil {
				log.Println(err)
			}
		}
	}()
}

func oldFunc(t string) string {
	return t
}
func newFunc(event *Event) string {
	return event.Data
}
