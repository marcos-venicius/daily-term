package idcluster

import "math/rand/v2"

const DefaultIdMaxSize = 9999

type IdCluster struct {
	history    map[int]bool // current in use id's
	randomizer *rand.Rand
	maxIdSize  int // max size of id
}

func CreateIdCluster() *IdCluster {
	r := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))

	return &IdCluster{
		history:    map[int]bool{},
		randomizer: r,
		maxIdSize:  DefaultIdMaxSize,
	}
}

func (cluster *IdCluster) SetCustomIdMaxSize(size int) {
	cluster.maxIdSize = size
}

// Generate a new unique id (inside this application and based on current existent values)
// Min 0, Max 9999 (default value ; can be modified by SetCustomIdMaxSize)
func (cluster *IdCluster) NewId() int {
	id := cluster.randomizer.IntN(cluster.maxIdSize)

	if _, ok := cluster.history[id]; ok {
		return cluster.NewId()
	}

	cluster.history[id] = true

	return id
}
