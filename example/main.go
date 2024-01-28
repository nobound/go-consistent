package main

import (
	"fmt"

	chash "github.com/nobound/go-consistent"
)

type Topics []string

func assignTopics(hash *chash.ConsistentHash, id string, topics Topics) Topics {
	var out Topics
	for _, tp := range topics {
		member := hash.Get(tp)
		if member == id {
			out = append(out, tp)
		}
	}
	return out
}

func main() {
	nodes := []string{"node1", "node2"}
	topics := []string{"topic1", "topic2", "topic3", "topic4", "topic5", "topic6"}
	replicas := 16

	showAssignment := func(hash *chash.ConsistentHash, nodes []string) {
		// nodes := hash.GetNodeNames()
		for _, node := range nodes {
			out := assignTopics(hash, node, topics)
			fmt.Printf("Topics assigned to %s\n", node)
			fmt.Printf("- number: %d\n", len(out))
			fmt.Printf("- topics: %v\n", out)
		}
	}

	config := chash.Config{
		ReplicationFactor: replicas,
	}
	hash := chash.NewWithNodes(nodes, config)
	showAssignment(hash, nodes)

	hash.Add("node3")
	nodes = append(nodes, "node3")
	showAssignment(hash, nodes)

	hash.Remove("node1")
	nodes = nodes[1:]
	showAssignment(hash, nodes)
}
