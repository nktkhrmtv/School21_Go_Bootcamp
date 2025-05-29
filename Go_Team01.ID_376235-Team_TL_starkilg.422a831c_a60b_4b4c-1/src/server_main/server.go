package main

import (
	"flag"
	"log"
	
    "team01/server"
)

func main() {
    host := flag.String("host", "127.0.0.1", "Host to bind")
    port := flag.String("port", "50051", "Port to bind")
    replicationFactor := flag.Int("rf", 2, "Replication factor")
    bootstrapHost := flag.String("bh", "", "Bootstrap host (optional)")
    bootstrapPort := flag.String("bp", "", "Bootstrap port (optional)")
    flag.Parse()

    self := server.Node{Host: *host, Port: *port}
    db := server.NewDatabase(self, *replicationFactor)

    var bootstrapNode server.Node
    if *bootstrapHost != "" && *bootstrapPort != "" {
        bootstrapNode = server.Node{Host: *bootstrapHost, Port: *bootstrapPort}
        log.Printf("Bootstrapping from %s", bootstrapNode)
    } else {
        log.Println("Starting as first node in cluster")
    }

    if bootstrapNode.String() != "" {
		server.JoinCluster(db, bootstrapNode)
	}

    go server.NodesHeartbeat(db)
    server.StartServer(db)
}

