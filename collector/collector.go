package collector

import (
	"flag"
)

var collectors []Collector

type Collector interface {
	Flags(flags *flag.FlagSet)
	Init() error
	Next() (string, func() (interface{}, error), error)
	Cleanup()
}

func Register(collector Collector) {
	collectors = append(collectors, collector)
}

func Collectors() []Collector {
	return collectors
}
