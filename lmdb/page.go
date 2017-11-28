package lmdb

import "unsafe"

type pageid uint64

type page struct {
	pgid     pageid
	flags    uint16
	count    uint16
	overflow uint32
	ptr      uintptr
}

type elemBranch struct {
	pos   uint32
	ksize uint32
}

type elemLeaf struct {
	pos   uint32
	ksize uint32
	vsize uint32
}

func (p *page) getMeta() *meta {
	return (*meta)(unsafe.Pointer(&p.ptr))
}
