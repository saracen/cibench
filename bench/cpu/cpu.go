package cpu

import (
	"flag"
	"hash/fnv"

	"github.com/saracen/cibench/bench"
)

const defaultIterations = 50000

var payload = []byte("the quick brown fox jumps over the lazy dog")

func init() {
	bench.Register(&cpu{})
}

type cpu struct {
	next       int
	iterations int
}

func (b *cpu) Flags(f *flag.FlagSet) {
	f.IntVar(&b.iterations, "cpu-iterations", defaultIterations, "cpu iterations per thread")
}

func (b *cpu) Init(threads int) error {
	return nil
}

func (b *cpu) Cleanup() {
}

func (b *cpu) Next() (string, func(threadIndex int) error, error) {
	defer func() {
		b.next++
	}()

	if b.next == 0 {
		return "cpu:fnv1a", b.fnv1a, nil
	}

	return "", nil, nil
}

func (b *cpu) fnv1a(threadIndex int) error {
	h := fnv.New64a()

	for i := 0; i < b.iterations; i++ {
		h.Write(payload)
		h.Sum64()
	}

	return nil
}
