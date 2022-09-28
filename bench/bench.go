package bench

import (
	"errors"
	"flag"
)

var benchers []Bencher

// ErrBenchEnded indicates the benchmark has ended, but isn't an error
var ErrBenchEnded = errors.New("bench ended")

type Bencher interface {
	Flags(flags *flag.FlagSet)
	Init(threads int) error
	Next() (string, func(threadIndex int) error, error)
	Cleanup()
}

func Register(bencher Bencher) {
	benchers = append(benchers, bencher)
}

func Benchers() []Bencher {
	return benchers
}
