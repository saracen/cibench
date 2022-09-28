package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/saracen/cibench/bench"
	"github.com/saracen/cibench/collector"

	_ "github.com/saracen/cibench/bench/cpu"
	_ "github.com/saracen/cibench/bench/fileio"
	_ "github.com/saracen/cibench/collector/internet"
)

type Info struct {
	Threads     int                `json:"threads"`
	Benchmarks  []BenchResult      `json:"benchmarks"`
	Collections []CollectionResult `json:"collection"`
}

type BenchResult struct {
	Test     string        `json:"test"`
	Total    int           `json:"total"`
	Duration time.Duration `json:"duration"`
}

type Result interface{}

type CollectionResult struct {
	Result
	Test string
}

var threads int

func main() {
	setup()

	var info Info
	err := run(&info)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	info.Threads = threads

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	enc.Encode(info)
}

func setup() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.IntVar(&threads, "threads", 1, "threads")

	for _, bencher := range bench.Benchers() {
		bencher.Flags(fs)
	}

	fs.Parse(os.Args[1:])
}

func run(info *Info) error {
	for _, bencher := range bench.Benchers() {
		err := bencher.Init(threads)
		if err != nil {
			return err
		}
		defer bencher.Cleanup()

		for {
			name, fn, err := bencher.Next()
			if err != nil {
				return err
			}
			if name == "" {
				break
			}

			total, duration := do(threads, fn)

			info.Benchmarks = append(info.Benchmarks, BenchResult{
				Test:     name,
				Total:    total,
				Duration: duration,
			})
		}
	}

	for _, collector := range collector.Collectors() {
		err := collector.Init()
		if err != nil {
			return err
		}

		for {
			name, fn, err := collector.Next()
			if err != nil {
				return err
			}
			if name == "" {
				break
			}

			result, err := fn()
			if err != nil {
				return err
			}

			info.Collections = append(info.Collections, CollectionResult{
				Test:   name,
				Result: result,
			})
		}
	}

	return nil
}
