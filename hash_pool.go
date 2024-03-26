package merkledag

import "hash"

type HashPool interface {
	Get() hash.Hash
}
