
# Go-Consistent

The `go-consistent` module serves as an implementation of the consistent hashing algorithm in the Go programming language. This algorithm is widely recognized and found applications in various domains. 

Consistent hashing offers an effective solution to the challenge of rehashing by providing a mechanism to distribute the load evenly across a distributed system. This approach not only mitigates the need for frequent rehashing but also ensures that the distribution of data remains balanced even as nodes are added or removed from the system.

For a comprehensive understanding of the consistent hashing algorithm and its practical implementations, additional insights can be obtained through a simple online search. (e.g., https://ably.com/blog/implementing-efficient-consistent-hashing)

## Usage

Import the module:
```go
import "github.com/nobound/go-consistent"
```

Use the module in your program:
```go
// Find all the topic assigned to a particular node
func assignTopics(hash *chash.ConsistentHash, node string, topics Topics) Topics {
	var out Topics
	for _, tp := range topics {
		member := hash.Get(tp)
		if member == node {
			out = append(out, tp)
		}
	}
	return out
}

// Assign multiple topics among the nodes
func main() {
	nodes := []string{"node1", "node2"}
	topics := []string{"topic1", "topic2", "topic3", "topic4", "topic5", "topic6"}
	replicas := 16

	showAssignment := func(hash *chash.ConsistentHash, nodes []string) {
		for _, node := range nodes {
			out := assignTopics(hash, node, topics)
			fmt.Printf("Topics assigned to %s\n", node)
			fmt.Printf("- number: %d\n", len(out))
			fmt.Printf("- topics: %v\n", out)
		}
	}

	config := Config{
		ReplicationFactor: replicas,
	}
	hash := NewWithNodes(nodes, config)
	showAssignment(hash, nodes)
}

```

## Distribution and Performance

The implementation is utilized to distribute 64,000 distinct topics across 10 nodes, taking into account different replica counts. This is used to assess both the implementation's effectiveness in distributing topics among the nodes and the efficiency of the calculation process.

| Number of replicas | Total Duration (ms) | Standard Deviation |
|               ---: |                ---: |               ---: | 
|                 50 |                 770 |            1141.21 |
|                100 |                 791 |             504.31 |
|                200 |                 809 |             450.69 |


