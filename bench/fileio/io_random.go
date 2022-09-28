package fileio

func (b *fileio) randomWrite(idx int) error {
	offset := b.rands[idx].Int63n(b.filesize - int64(len(b.blocks[idx])+1))

	_, err := b.files[idx].WriteAt(b.blocks[idx], int64(offset))

	return err
}

func (b *fileio) randomRead(idx int) error {
	offset := b.rands[idx].Int63n(b.filesize - int64(len(b.blocks[idx])+1))

	_, err := b.files[idx].ReadAt(b.blocks[idx], int64(offset))

	return err
}
