package fileio

import "io"

func (b *fileio) sequentialWrite(idx int) error {
	n, err := b.files[idx].Write(b.blocks[idx])

	b.offsets[idx] += int64(n)
	if b.offsets[idx] >= b.filesize {
		if err := b.reset(idx); err != nil {
			return err
		}
	}

	return err
}

func (b *fileio) sequentialRead(idx int) error {
	n, err := b.files[idx].Read(b.blocks[idx])

	b.offsets[idx] += int64(n)
	if b.offsets[idx] >= b.filesize {
		if err := b.reset(idx); err != nil {
			return err
		}
	}

	return err
}

func (b *fileio) reset(idx int) error {
	_, err := b.files[idx].Seek(0, io.SeekStart)
	b.offsets[idx] = 0

	return err
}
