package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	
	"team01/warehouse_rpc"
)

func (c *Client) updateClusterInfo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := c.client.Heartbeat(ctx, &proto.HeartbeatRequest{
		Node: &proto.NodeInfo{Host: "client", Port: "0"},
	})
	if err != nil {
		return fmt.Errorf("heartbeat failed: %v", err)
	}

	c.nodes = nil
	for _, node := range resp.Nodes {
		c.nodes = append(c.nodes, fmt.Sprintf("%s:%s", node.Host, node.Port))
	}
	c.replicationFactor = int(resp.ReplicationFactor)
	return nil
}

func (c *Client) monitorNodes() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		if err := c.updateClusterInfo(); err != nil {
			log.Printf("Failed to update cluster info: %v", err)
		}

		var aliveNodes []string
		currentNodeAlive := false

		for _, node := range c.nodes {
			if ok := c.checkNodeAlive(node); ok {
				aliveNodes = append(aliveNodes, node)
				if node == c.currentNode {
					currentNodeAlive = true
				}
			}
		}

		c.nodes = aliveNodes
		if !currentNodeAlive && len(c.nodes) > 0 {
			c.currentNode = c.nodes[0]
			c.updateConnection()
			fmt.Printf("\nReconnected to node %s\n", c.currentNode)
			c.PrintNodes()
		}

		if len(c.nodes) < c.replicationFactor {
			fmt.Printf("\nWARNING: cluster size (%d) < replication factor (%d)!\n",
				len(c.nodes), c.replicationFactor)
		}
		c.mu.Unlock()
	}
}

func (c *Client) checkNodeAlive(node string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		node,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return false
	}
	defer conn.Close()

	client := proto.NewWarehouseClient(conn)
	_, err = client.Heartbeat(ctx, &proto.HeartbeatRequest{
		Node: &proto.NodeInfo{Host: "client", Port: "0"},
	})
	return err == nil
}

func (c *Client) PrintNodes() {
	fmt.Println("Known nodes:")
	for _, node := range c.nodes {
		fmt.Println(node)
	}
}

func (c *Client) CurrentNode() any {
	return c.currentNode
}