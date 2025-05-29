package server

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"team01/warehouse_rpc"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Database struct {
	proto.UnimplementedWarehouseServer
	data              map[string]string
	mu                sync.RWMutex
	nodes             []Node
	replicationFactor int
	self              Node
}

func NewDatabase(self Node, replicationFactor int) *Database {
	return &Database{
		data:              make(map[string]string),
		nodes:             []Node{self},
		replicationFactor: replicationFactor,
		self:              self,
	}
}

func (db *Database) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	if _, err := uuid.Parse(req.Key); err != nil {
		return nil, errors.New("key is not a valid UUID4")
	}

	db.mu.RLock()
	val, ok := db.data[req.Key]
	db.mu.RUnlock()

	return &proto.GetResponse{
		Value: val,
		Found: ok,
	}, nil
}

func (db *Database) Set(ctx context.Context, req *proto.SetRequest) (*proto.SetResponse, error) {
	if _, err := uuid.Parse(req.Key); err != nil {
        return nil, errors.New("key is not a valid UUID4")
    }

    // Проверяем, является ли запрос репликацией
    if req.IsReplica {
        // Просто записываем данные без репликации
        db.mu.Lock()
        db.data[req.Key] = req.Value
        db.mu.Unlock()
        return &proto.SetResponse{Replicas: 1}, nil
    }

    // 1. Локальная запись
    db.mu.Lock()
    db.data[req.Key] = req.Value
    db.mu.Unlock()

    // 2. Репликация на другие узлы
    replicas := db.getReplicaNodes(req.Key)
    successReplicas := 1 

    for _, node := range replicas {
        if node == db.self {
            continue
        }
		grpc.NewServer()
		conn, err := grpc.NewClient(node.String(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
        if err != nil {
            log.Printf("Replication failed to %s: %v", node, err)
            continue
        }

        client := proto.NewWarehouseClient(conn)
        nodeCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
        
        // Отправляем запрос с флагом IsReplica = true
        _, err = client.Set(nodeCtx, &proto.SetRequest{
            Key:       req.Key,
            Value:     req.Value,
            IsReplica: true, 
        })
        conn.Close()
        cancel()

        if err == nil {
            successReplicas++
        }
    }

    return &proto.SetResponse{
        Replicas: int32(successReplicas),
    }, nil
}

func (db *Database) Delete(ctx context.Context, req *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	if _, err := uuid.Parse(req.Key); err != nil {
		return nil, errors.New("key is not a valid UUID4")
	}

	db.mu.Lock()
	delete(db.data, req.Key)
	db.mu.Unlock()

	replicas := db.getReplicaNodes(req.Key)
	return &proto.DeleteResponse{
		Replicas: int32(len(replicas)),
	}, nil
}

func (db *Database) Heartbeat(ctx context.Context, req *proto.HeartbeatRequest) (*proto.HeartbeatResponse, error) {
	nodes := db.getNodes()
	var nodeInfos []*proto.NodeInfo
	for _, n := range nodes {
		nodeInfos = append(nodeInfos, &proto.NodeInfo{
			Host: n.Host,
			Port: n.Port,
		})
	}

	return &proto.HeartbeatResponse{
		Nodes:             nodeInfos,
		ReplicationFactor: int32(db.replicationFactor),
		ReplicationWarning: len(nodes) < db.replicationFactor,
	}, nil
}

func (db *Database) Join(ctx context.Context, req *proto.JoinRequest) (*proto.JoinResponse, error) {
	if int(req.ReplicationFactor) != db.replicationFactor {
		return &proto.JoinResponse{Accepted: false}, nil
	}

	newNode := Node{Host: req.NewNode.Host, Port: req.NewNode.Port}
	db.addNode(newNode)

	for _, node := range db.getNodes() {
		if node == db.self || node == newNode {
			continue
		}
		go db.notifyNodeAboutNewNode(node, newNode)
	}

	return &proto.JoinResponse{
		Nodes:    db.getNodeInfos(),
		Accepted: true,
	}, nil
}

func (db *Database) SyncData(req *proto.SyncRequest, stream proto.Warehouse_SyncDataServer) error {
    db.mu.RLock()
    defer db.mu.RUnlock()

    for key, value := range db.data {
        if err := stream.Send(&proto.DataEntry{
            Key:   key,
            Value: value,
        }); err != nil {
            return err
        }
    }
    return nil
}

