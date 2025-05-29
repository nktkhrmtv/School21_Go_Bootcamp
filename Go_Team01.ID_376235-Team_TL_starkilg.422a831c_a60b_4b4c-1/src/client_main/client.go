package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	
	"team01/client"
)

func main() {
	host := flag.String("H", "127.0.0.1", "Host to connect")
	port := flag.String("P", "8765", "Port to connect")
	flag.Parse()

	c, err := client.NewClient(*host, *port)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer c.Close()

	fmt.Printf("Connected to node %s\n", c.CurrentNode())
	c.PrintNodes()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		cmd := scanner.Text()
		parts := strings.Fields(cmd)
		if len(parts) == 0 {
			fmt.Print("> ")
			continue
		}

		switch parts[0] {
		case "GET":
			handleGet(c, parts)
		case "SET":
			handleSet(c, parts)
		case "DELETE":
			handleDelete(c, parts)
		case "NODES":
			c.PrintNodes()
		default:
			fmt.Println("Unknown command. Available: GET, SET, DELETE, NODES")
		}
		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func handleGet(c *client.Client, parts []string) {
	if len(parts) != 2 {
		fmt.Println("Usage: GET <key>")
		return
	}
	val, err := c.Get(parts[1])
	if err != nil {
		if err.Error() == "not found" {
			fmt.Println("Not found")
		} else {
			fmt.Printf("Error: %v\n", err)
		}
	} else {
		fmt.Println(val)
	}
}

func handleSet(c *client.Client, parts []string) {
	if len(parts) < 3 {
		fmt.Println("Usage: SET <key> <value>")
		return
	}
	replicas, err := c.Set(parts[1], strings.Join(parts[2:], " "))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Created (%d replicas)\n", replicas)
	}
}

func handleDelete(c *client.Client, parts []string) {
	if len(parts) != 2 {
		fmt.Println("Usage: DELETE <key>")
		return
	}
	replicas, err := c.Delete(parts[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Deleted (%d replicas)\n", replicas)
	}
}