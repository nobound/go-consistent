package go_consistent_hash

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"sync"
)

var show bool = false

type Bignum struct {
	*big.Int
}

func NewBignum(n string) *Bignum {
	x1, ok := new(big.Int).SetString(n, 0)
	if !ok {
		fmt.Println("fail to convert")
		return nil
	}
	return &Bignum{x1}
}

type SortedKeys []*Bignum

func (x SortedKeys) Len() int {
	return len(x)
}

func (x SortedKeys) Less(i, j int) bool {
	rv := x[i].Cmp(x[j].Int) < 0
	if show {
		if rv {
			fmt.Println(x[i], " < ")
			fmt.Println(x[j])
		} else {
			fmt.Println(x[i], " >= ")
			fmt.Println(x[j])
		}
	}
	return rv
}

func (x SortedKeys) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func md5_sum(s []byte) *Bignum {
	out := md5.Sum(s)
	// fmt.Println("byte array:", out)
	hexstr := hex.EncodeToString(out[:])
	// fmt.Println("encoded hex string: ", hexstr)
	num, _ := new(big.Int).SetString(hexstr, 16)
	return &Bignum{num}
}

type Config struct {
	ReplicationFactor int
}

type ConsistentHashing struct {
	config         Config
	sortedHashKeys SortedKeys
	ring           map[string]string
	dataSet        map[string]bool
	mu             sync.Mutex
}

func New(members []string, config Config) *ConsistentHashing {
	c := &ConsistentHashing{
		config:  config,
		ring:    make(map[string]string),
		dataSet: make(map[string]bool),
	}
	for _, m := range members {
		c.Add(m)
	}
	return c
}

// Get a nearest object name from input object in consistent hashing ring
func (c *ConsistentHashing) Get(key string) string {
	index := c.searchRingIndex(key)
	skey := c.sortedHashKeys[index]
	s := skey.String()
	node, found := c.ring[s]
	if !found {
		fmt.Println("cannot find value for key ", skey)
	}
	return node
}

func (c *ConsistentHashing) Add(nodename string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	names := c.getNodeKeys(nodename)
	for vname, hkey := range names {
		s := hkey.String()
		c.ring[s] = nodename
		c.dataSet[vname] = true
	}
	c.updateSortHashKeys()
}

func (c *ConsistentHashing) Remove(nodename string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	names := c.getNodeKeys(nodename)
	for vname, hkey := range names {
		delete(c.ring, hkey.String())
		delete(c.dataSet, vname)
	}
	c.updateSortHashKeys()
}

// TBD
func (c *ConsistentHashing) GetNodeNames() []string {
	var out []string
	for k, _ := range c.dataSet {
		out = append(out, k)
	}
	return out
}

func (c *ConsistentHashing) getNodeKeys(nodename string) map[string]*Bignum {
	out := make(map[string]*Bignum)
	for i := 0; i < c.config.ReplicationFactor; i++ {
		s := fmt.Sprintf("%s:%d", nodename, i)
		h := c.hasKey(s)
		out[s] = h
	}
	return out
}

// The node replica with a hash value nearest but not less than that of the given
// name is returned.   If the hash of the given name is greater than the greatest
// hash, returns the lowest hashed node.

func (c *ConsistentHashing) searchRingIndex(obj string) int {
	count := len(c.sortedHashKeys)
	targetKey := c.hasKey(obj)

	// big.num compare function x.Cmp(y)
	// -1 if x <  y
	//  0 if x == y
	// +1 if x >  y
	fn := func(i int) bool {
		x := c.sortedHashKeys[i]
		y := targetKey
		rv := x.Cmp(y.Int) > 0
		// debug
		if show {
			fmt.Println(i)
			if rv {
				// ture when x >= y
				fmt.Println(x, " >= ")
				fmt.Println(y)
			} else {
				// false when x < y
				fmt.Println(x, " < ")
				fmt.Println(y)
			}
		}
		return rv
	}

	targetIndex := sort.Search(count, fn)
	if targetIndex >= count {
		targetIndex = 0
	}
	return targetIndex
}

func (c *ConsistentHashing) updateSortHashKeys() {
	c.sortedHashKeys = nil
	for nodename, _ := range c.dataSet {
		key := c.hasKey(nodename)
		c.sortedHashKeys = append(c.sortedHashKeys, key)
	}
	sort.Sort(c.sortedHashKeys)
}

func (c *ConsistentHashing) hasKey(obj string) *Bignum {
	return md5_sum([]byte(obj))
}
