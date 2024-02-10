package chash

import (
	"fmt"
	"testing"
)

var numberOfNodes = 10
var numberOfTopics = 64000

type Node struct {
	id     string
	topics []string
}

func createNodes() []*Node {
	count := numberOfNodes
	var nodes []*Node
	for i := 0; i < count; i++ {
		s := fmt.Sprintf("node%d", i)
		n := Node{
			id: s,
		}
		nodes = append(nodes, &n)
	}
	return nodes
}

type Topics []string

func createTopics() Topics {
	count := numberOfTopics
	topics := make(Topics, count)
	for i := 0; i < count; i++ {
		s := fmt.Sprintf("topic%d", i)
		topics[i] = s
	}
	return topics
}

func assignTopics(hash *ConsistentHash, node *Node, topics Topics) {
	for _, tp := range topics {
		id := hash.Get(tp)
		if id == node.id {
			node.topics = append(node.topics, tp)
		}
	}
}

type set map[string]struct{}

// Track the unique values encountered across all arrays using a hash set.
// Use map for the value is encountered more than once.
func checkTopics(nodes []*Node) (map[string]int, set) {
	occurrences := make(map[string]int)
	uniqueValues := make(set)
	for _, n := range nodes {
		for _, val := range n.topics {
			if _, ok := uniqueValues[val]; !ok {
				uniqueValues[val] = struct{}{}
			} else {
				occurrences[val]++
			}
		}
	}
	return occurrences, uniqueValues
}

func TestConsistentHash(t *testing.T) {
	nodes := createNodes()
	topics := createTopics()

	fmt.Printf("Number of nodes: %d, number of topics: %d\n", len(nodes), len(topics))

	replicas := 200
	config := Config{
		ReplicationFactor: replicas,
	}
	hash := New(config)

	// add nodes to the hash ring
	for _, n := range nodes {
		hash.Add(n.id)
	}

	// assign topics to each node
	for _, n := range nodes {
		assignTopics(hash, n, topics)
	}

	occurrences, unique := checkTopics(nodes)
	for val, count := range occurrences {
		if count > 1 {
			t.Errorf("Value %s is duplicated %d times\n", val, count)
		}
	}

	if len(topics) != len(unique) {
		t.Errorf("Before distribution: %d != After distribution: %d\n",
			len(topics), len(unique))
	}
}
