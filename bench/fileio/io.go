package fileio

import (
	cryptorand "crypto/rand"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/saracen/cibench/bench"
)

const (
	nextRandomWrite = iota
	nextRandomRead
	nextSequentialWrite
	nextSequentialRead

	defaultBlocksize = 1024 * 1024
	defaultFilesize  = 1024 * 1024 * 1024 * 2
)

func init() {
	bench.Register(&fileio{})
}

type fileio struct {
	next int

	files   []*os.File
	blocks  [][]byte
	rands   []*rand.Rand
	offsets []int64

	filesize  int64
	blocksize int64
}

func (b *fileio) Flags(f *flag.FlagSet) {
	f.Int64Var(&b.filesize, "io-file-size", defaultFilesize, "io file size")
	f.Int64Var(&b.blocksize, "io-block-size", defaultBlocksize, "io block size")
}

func (b *fileio) Init(threads int) error {
	if b.filesize < b.blocksize {
		return fmt.Errorf("filesize cannot be less than blocksize")
	}

	b.offsets = make([]int64, threads)
	b.blocks = make([][]byte, threads)

	for i := 0; i < threads; i++ {
		f, err := os.OpenFile(fmt.Sprintf("cibench_data_%d", i), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return fmt.Errorf("creating data file: %w", err)
		}
		if _, err := io.CopyN(f, cryptorand.Reader, b.filesize); err != nil {
			return fmt.Errorf("writing random data: %w", err)
		}
		if err := f.Sync(); err != nil {
			return fmt.Errorf("syncing random data: %w", err)
		}

		b.blocks[i] = make([]byte, b.blocksize)
		if _, err := cryptorand.Read(b.blocks[i]); err != nil {
			return fmt.Errorf("writing random data to block: %w", err)
		}

		b.files = append(b.files, f)
		b.rands = append(b.rands, rand.New(rand.NewSource(time.Now().UnixNano()+int64(i))))
	}

	return nil
}

func (b *fileio) Cleanup() {
	for _, f := range b.files {
		f.Close()
		os.Remove(f.Name())
	}
}

func (b *fileio) Next() (string, func(threadIndex int) error, error) {
	defer func() {
		b.next++
	}()

	switch b.next {
	case nextRandomWrite:
		return "io:random-write", b.randomWrite, nil
	case nextRandomRead:
		return "io:random-read", b.randomRead, nil
	case nextSequentialWrite:
		for i := range b.files {
			if err := b.reset(i); err != nil {
				return "", nil, fmt.Errorf("io:sequential-read reset error: %w", err)
			}
		}
		return "io:sequential-read", b.sequentialRead, nil
	case nextSequentialRead:
		for i := range b.files {
			if err := b.reset(i); err != nil {
				return "", nil, fmt.Errorf("io:sequential-write reset error: %w", err)
			}
		}
		return "io:sequential-write", b.sequentialWrite, nil
	}

	return "", nil, nil
}
