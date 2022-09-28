package internet

import (
	"context"
	"flag"
	"time"

	"github.com/m-lab/ndt7-client-go"

	"github.com/saracen/cibench/collector"
)

func init() {
	collector.Register(&internet{})
}

type internet struct {
	next   int
	client *ndt7.Client
}

func (b *internet) Flags(f *flag.FlagSet) {
}

func (b *internet) Init() error {
	b.client = ndt7.NewClient("ndt7-client-go-example", "0.1.0")

	return nil
}

func (b *internet) Cleanup() {
}

func (b *internet) Next() (string, func() (interface{}, error), error) {
	defer func() {
		b.next++
	}()

	if b.next == 0 {
		return "internet:upload", b.upload, nil
	}

	if b.next == 1 {
		return "internet:download", b.download, nil
	}

	return "", nil, nil
}

func (b *internet) download() (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ch, err := b.client.StartDownload(ctx)
	if err != nil {
		return nil, err
	}

	// drain
	for range ch {
	}

	return b.client.Results(), nil
}

func (b *internet) upload() (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ch, err := b.client.StartUpload(ctx)
	if err != nil {
		return nil, err
	}

	// drain
	for range ch {
	}

	return b.client.Results(), nil
}
