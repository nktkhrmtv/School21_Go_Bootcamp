package client

import (
	"context"
	"errors"
	"team01/common"
	"team01/warehouse_rpc"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (c *Client) Get(key string) (string, error) {
	if _, err := uuid.Parse(key); err != nil {
		return "", errors.New("key is not a valid UUID4")
	}

	for _, node := range c.getReplicaNodes(key) {
		value, err := c.getFromNode(node, key)
		if err == nil {
			return value, nil
		}
	}
	return "", errors.New("all nodes unavailable")
}

func (c *Client) getFromNode(node, key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		node,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	resp, err := proto.NewWarehouseClient(conn).Get(ctx, &proto.GetRequest{Key: key})
	if err != nil {
		return "", err
	}
	if !resp.Found {
		return "", errors.New("not found")
	}
	return resp.Value, nil
}

func (c *Client) Set(key, value string) (int, error) {
	if _, err := uuid.Parse(key); err != nil {
		return 0, errors.New("key is not a valid UUID4")
	}

	for _, node := range c.nodes {
		replicas, err := c.setToNode(node, key, value)
		if err == nil {
			return replicas, nil
		}
	}
	return 0, errors.New("all nodes unavailable")
}

func (c *Client) setToNode(node, key, value string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		node,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	resp, err := proto.NewWarehouseClient(conn).Set(ctx, &proto.SetRequest{
		Key:   key,
		Value: value,
		IsReplica: false,
	})
	if err != nil {
		return 0, err
	}
	return int(resp.Replicas), nil
}

func (c *Client) Delete(key string) (int, error) {
	if _, err := uuid.Parse(key); err != nil {
		return 0, errors.New("key is not a valid UUID4")
	}

	for _, node := range c.nodes {
		replicas, err := c.deleteFromNode(node, key)
		if err == nil {
			return replicas, nil
		}
	}
	return 0, errors.New("all nodes unavailable")
}

func (c *Client) deleteFromNode(node, key string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		node,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	resp, err := proto.NewWarehouseClient(conn).Delete(ctx, &proto.DeleteRequest{Key: key})
	if err != nil {
		return 0, err
	}
	return int(resp.Replicas), nil
}

func (c *Client) getReplicaNodes(key string) []string {
	if len(c.nodes) <= c.replicationFactor {
		return c.nodes
	}

	hash := common.FnvHash(key)
	startIdx := hash % uint32(len(c.nodes))
	
	var replicas []string
	for i := 0; i < c.replicationFactor; i++ {
		idx := (startIdx + uint32(i)) % uint32(len(c.nodes))
		replicas = append(replicas, c.nodes[idx])
	}
	return replicas
}