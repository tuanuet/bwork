package bwork

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/stathat/consistent"
	pb "github.com/tuanuet/bwork/proto"
)

// ShardManager manages records and nodes with gRPC streaming
type ShardManager struct {
	hash     *consistent.Consistent
	records  map[string]bool
	nodes    map[string]bool
	clients  map[string]chan *pb.RebalanceEvent
	lastSeen map[string]time.Time
	mu       sync.Mutex
	pb.UnimplementedShardServiceServer
}

// NewShardManager initializes the sharding system
func NewShardManager() *ShardManager {
	sm := &ShardManager{
		hash:     consistent.New(),
		records:  make(map[string]bool),
		nodes:    make(map[string]bool),
		clients:  make(map[string]chan *pb.RebalanceEvent),
		lastSeen: make(map[string]time.Time),
	}
	go sm.checkNodeAvailability()
	return sm
}

// RegisterNode registers a new node
func (sm *ShardManager) RegisterNode(ctx context.Context, req *pb.NodeRequest) (*pb.RegisterResponse, error) {
	sm.mu.Lock()

	nodeID := req.NodeId
	if sm.nodes[nodeID] {
		sm.mu.Unlock()
		return &pb.RegisterResponse{Success: false, Message: "Node already registered"}, nil
	}

	sm.hash.Add(nodeID)
	sm.nodes[nodeID] = true
	sm.clients[nodeID] = make(chan *pb.RebalanceEvent, 100)
	sm.lastSeen[nodeID] = time.Now()
	fmt.Printf("Registered node %s\n", nodeID)
	sm.mu.Unlock()
	sm.rebalanceRecords()
	return &pb.RegisterResponse{Success: true, Message: "Node registered successfully"}, nil
}

// Ping updates the node's last-seen timestamp
func (sm *ShardManager) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	nodeID := req.NodeId
	if !sm.nodes[nodeID] {
		return &pb.PingResponse{Success: false}, nil
	}

	sm.lastSeen[nodeID] = time.Now()
	return &pb.PingResponse{Success: true}, nil
}

// checkNodeAvailability removes unresponsive nodes
func (sm *ShardManager) checkNodeAvailability() {
	for {
		time.Sleep(1 * time.Second)
		//sm.mu.Lock()
		now := time.Now()
		for nodeID, lastSeen := range sm.lastSeen {
			if now.Sub(lastSeen) > 2*time.Second {
				fmt.Printf("Node %s unresponsive, removing\n", nodeID)
				sm.hash.Remove(nodeID)
				delete(sm.nodes, nodeID)
				close(sm.clients[nodeID])
				delete(sm.clients, nodeID)
				delete(sm.lastSeen, nodeID)
				sm.rebalanceRecords()
			}
		}
		//sm.mu.Unlock()
	}
}

// AddRecord adds a record and notifies the assigned node
func (sm *ShardManager) AddRecord(recordID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.records[recordID] = true
	fmt.Printf("Added record %s\n", recordID)
	sm.notifyNode(recordID, false)
}

// RemoveRecord removes a record and notifies the assigned node
func (sm *ShardManager) RemoveRecord(ctx context.Context, req *pb.RemoveRecordRequest) (*pb.RemoveRecordResponse, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	recordID := req.RecordId
	if _, exists := sm.records[recordID]; !exists {
		return &pb.RemoveRecordResponse{Success: false, Message: "Event not found"}, nil
	}

	delete(sm.records, recordID)
	fmt.Printf("Removed record %s\n", recordID)
	sm.notifyNode(recordID, true)

	return &pb.RemoveRecordResponse{Success: true, Message: "Event removed successfully"}, nil
}

// rebalanceRecords notifies old and new nodes of record movements
func (sm *ShardManager) rebalanceRecords() {
	fmt.Println("Rebalancing records across all nodes")
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Notify all nodes of the rebalancing event
	for _, ch := range sm.clients {
		ch <- &pb.RebalanceEvent{}
	}
}

// notifyNode sends a record update to the assigned node and tracks it
func (sm *ShardManager) notifyNode(recordID string, isDeleted bool) {
	node, err := sm.hash.Get(recordID)
	if err != nil {
		log.Printf("Error assigning record %s: %v", recordID, err)
		return
	}
	if ch, ok := sm.clients[node]; ok {
		ch <- &pb.RebalanceEvent{}
	}
}

// GetRecordsForNode implements the unary RPC
func (sm *ShardManager) GetRecordsForNode(ctx context.Context, req *pb.NodeRequest) (*pb.RecordsResponse, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	nodeID := req.NodeId
	if !sm.nodes[nodeID] {
		return nil, fmt.Errorf("node %s not found", nodeID)
	}

	var records []*pb.Record
	for recordID, _ := range sm.records {
		assignedNode, err := sm.hash.Get(recordID)
		if err != nil {
			continue
		}
		if assignedNode == nodeID {
			records = append(records, &pb.Record{RecordId: recordID})
		}
	}
	return &pb.RecordsResponse{Records: records}, nil
}

// SubscribeRebalance streams updates to a node
func (sm *ShardManager) SubscribeRebalance(req *pb.NodeRequest, stream pb.ShardService_SubscribeRebalanceServer) error {
	nodeID := req.NodeId
	sm.mu.Lock()
	if !sm.nodes[nodeID] {
		sm.mu.Unlock()
		return fmt.Errorf("node %s not registered", nodeID)
	}
	ch := sm.clients[nodeID]
	sm.mu.Unlock()

	for update := range ch {
		if err := stream.Send(update); err != nil {
			return err
		}
	}
	return nil
}
