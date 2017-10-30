package lmdb

import "unsafe"

const (
	_PFLAG_BRANCH   = 0x01
	_PFLAG_LEAF     = 0x02
	_PFLAG_META     = 0x04
	_PFLAG_FREELIST = 0x08

	_PGID_META     = 0
	_PGID_FREELIST = 1
	_PGID_LEAF     = 2
)

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
