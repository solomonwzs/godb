package lmdb

type Bucket struct {
}

func (b *Bucket) NewCursor() *Cursor {
	return &Cursor{
		bucket: b,
	}
}
