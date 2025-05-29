package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"sort"
	"time"

	"team01/warehouse_rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Node struct {
	Host string
	Port string
}

func (n Node) String() string {
	return fmt.Sprintf("%s:%s", n.Host, n.Port)
}

func StartServer(db *Database) {
    lis, err := net.Listen("tcp", db.self.String())
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    proto.RegisterWarehouseServer(grpcServer, db)
    log.Printf("Server started at %s", db.self.String())
    
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

func (db *Database) addNode(node Node) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for _, n := range db.nodes {
		if n == node {
			return
		}
	}
	db.nodes = append(db.nodes, node)
	sort.Slice(db.nodes, func(i, j int) bool {
		return db.nodes[i].String() < db.nodes[j].String()
	})
}

func (db *Database) removeNode(node Node) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for i, n := range db.nodes {
		if n == node {
			db.nodes = append(db.nodes[:i], db.nodes[i+1:]...)
			break
		}
	}
}

func (db *Database) getNodes() []Node {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.nodes
}

func (db *Database) getNodeInfos() []*proto.NodeInfo {
	nodes := db.getNodes()
	var nodeInfos []*proto.NodeInfo
	for _, n := range nodes {
		nodeInfos = append(nodeInfos, &proto.NodeInfo{
			Host: n.Host,
			Port: n.Port,
		})
	}
	return nodeInfos
}

func NodesHeartbeat(db *Database) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop() 

    for range ticker.C {
        nodes := db.getNodes()

        for _, node := range nodes {
            if node == db.self {
                continue
            }

            // Для каждого узла запускаем проверку в отдельной горутине
            go func(node Node) {
				conn, err := grpc.NewClient(node.String(),
					grpc.WithTransportCredentials(insecure.NewCredentials()),
				)

                if err != nil {
                    // Если соединение не установлено, считаем узел недоступным
                    log.Printf("Node %s is down, removing from cluster", node)
                    db.removeNode(node) // Удаляем узел из кластера
                    return
                }
                defer conn.Close() 

                client := proto.NewWarehouseClient(conn)
                ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
                defer cancel()

                _, err = client.Heartbeat(ctx, &proto.HeartbeatRequest{
                    Node: &proto.NodeInfo{
                        Host: db.self.Host,
                        Port: db.self.Port,
                    },
                })
                if err != nil {
                    // Если запрос не прошел, считаем узел недоступным
                    log.Printf("Node %s is down, removing from cluster", node)
                    db.removeNode(node) // Удаляем узел из кластера
                }
            }(node) // node как параметр для избежания замыкания (работа с одним и тем же node)
        }
    }
}
