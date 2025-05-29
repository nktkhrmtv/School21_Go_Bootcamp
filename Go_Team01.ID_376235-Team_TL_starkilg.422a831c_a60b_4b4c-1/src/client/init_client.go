package client

import (
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"team01/warehouse_rpc"
)

type Client struct {
	conn              *grpc.ClientConn
	client            proto.WarehouseClient
	currentNode       string
	nodes             []string
	replicationFactor int
	mu                sync.Mutex
}

func NewClient(host, port string) (*Client, error) {
	node := fmt.Sprintf("%s:%s", host, port)

	conn, err := grpc.NewClient(
		node,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to node %s: %v", node, err)
	}

	c := &Client{
		conn:        conn,
		client:      proto.NewWarehouseClient(conn),
		currentNode: node,
		nodes:       []string{node},
	}

	if err := c.updateClusterInfo(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to get cluster info: %v", err)
	}

	go c.monitorNodes()
	return c, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) updateConnection() {
	if c.conn != nil {
		c.conn.Close()
	}

	c.conn, _ = grpc.NewClient(
		c.currentNode,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	c.client = proto.NewWarehouseClient(c.conn)
}