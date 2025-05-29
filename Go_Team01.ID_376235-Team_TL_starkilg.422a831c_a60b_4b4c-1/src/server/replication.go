package server

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	"team01/common"
	"team01/warehouse_rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (db *Database) getReplicaNodes(key string) []Node {
	nodes := db.getNodes()
	if len(nodes) <= db.replicationFactor {
		return nodes
	}

	hash := common.FnvHash(key)
	startIdx := hash % uint32(len(nodes))
	
	var replicas []Node
	for i := 0; i < db.replicationFactor; i++ {
		idx := (startIdx + uint32(i)) % uint32(len(nodes))
		replicas = append(replicas, nodes[idx])
	}
	return replicas
}

func (db *Database) notifyNodeAboutNewNode(node, newNode Node) {
	conn, err := grpc.NewClient(node.String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Printf("Failed to connect to node %s: %v", node, err)
		return
	}
	defer conn.Close()

	client := proto.NewWarehouseClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err = client.Join(ctx, &proto.JoinRequest{
		NewNode: &proto.NodeInfo{
			Host: newNode.Host,
			Port: newNode.Port,
		},
		ReplicationFactor: int32(db.replicationFactor),
	})
	if err != nil {
		log.Printf("Failed to notify node %s about new node: %v", node, err)
	}
}

func JoinCluster(db *Database, bootstrapNode Node) {
    if bootstrapNode.Host == "" || bootstrapNode.Port == "" {
        return
    }

    conn, err := grpc.NewClient(
        bootstrapNode.String(),
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
            return net.DialTimeout("tcp", addr, 2*time.Second)
        }),
    )
    if err != nil {
        log.Printf("Failed to connect to bootstrap node %s: %v", bootstrapNode, err)
        return
    }
    defer conn.Close()

    client := proto.NewWarehouseClient(conn)
    
    // Таймаут для Join-запроса 
    joinCtx, joinCancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer joinCancel()

    // 1. Присоединяемся к кластеру
    resp, err := client.Join(joinCtx, &proto.JoinRequest{
        NewNode: &proto.NodeInfo{
            Host: db.self.Host,
            Port: db.self.Port,
        },
        ReplicationFactor: int32(db.replicationFactor),
    })
    if err != nil {
        log.Printf("Failed to join cluster: %v", err)
        return
    }

    if !resp.Accepted {
        log.Println("Failed to join cluster: replication factor mismatch")
        return
    }

    // 2. Получаем список узлов
    for _, node := range resp.Nodes {
        db.addNode(Node{Host: node.Host, Port: node.Port})
    }

    // 3. Синхронизируем данные с лидером
    syncCtx, syncCancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer syncCancel()

    syncStream, err := client.SyncData(syncCtx, &proto.SyncRequest{})
    if err != nil {
        log.Printf("Failed to start data sync: %v", err)
        return
    }

    db.mu.Lock()
    defer db.mu.Unlock()

    for {
        entry, err := syncStream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Printf("Data sync error: %v", err)
            break
        }
        db.data[entry.Key] = entry.Value
    }

    log.Printf("Data sync completed, received %d entries", len(db.data))
}

