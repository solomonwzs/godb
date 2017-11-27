package lmdb

import "unsafe"

type pageid uint64

type page struct {
	id       pageid
	flags    uint16
	count    uint16
	overflow uint32
	ptr      uintptr
}

func (p *page) getMeta() *meta {
	return (*meta)(unsafe.Pointer(&p.ptr))
}
