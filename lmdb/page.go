package lmdb

import (
	"unsafe"
)

const (
	_PAGE_FLAG_UNKNOWN  = 0x00
	_PAGE_FLAG_BRANCH   = 0x01
	_PAGE_FLAG_LEAF     = 0x02
	_PAGE_FLAG_META     = 0x04
	_PAGE_FLAG_FREELIST = 0x08

	_LEAF_ELEM_SIZE   = int(unsafe.Sizeof(leafElem{}))
	_BRANCH_ELEM_SIZE = int(unsafe.Sizeof(branchElem{}))

	_MAX_ELEMENT_COUNT = 0xffff
)

type pageid uint64

type page struct {
	pgid     pageid
	flags    uint16
	count    uint16
	overflow uint32
	ptr      uintptr
}

type branchElem struct {
	pgid  pageid
	pos   uint32
	ksize uint32
}

type leafElem struct {
	pos   uint32
	ksize uint32
	vsize uint32
}

func (p *page) getMeta() *meta {
	if p.flags != _PAGE_FLAG_META {
		return nil
	}
	return (*meta)(unsafe.Pointer(&p.ptr))
}

func (p *page) getLeafElems() []leafElem {
	if p.count == 0 {
		return nil
	}
	return (*[_MAX_ELEMENT_COUNT]leafElem)(
		unsafe.Pointer(&p.ptr))[:p.count]
}

func (p *page) getBranchElems() []branchElem {
	if p.count == 0 {
		return nil
	}
	return (*[_MAX_ELEMENT_COUNT]branchElem)(
		unsafe.Pointer(&p.ptr))[:p.count]
}

func (e *leafElem) key() []byte {
	buf := (*[_MAX_ALLOC_SIZE]byte)(unsafe.Pointer(e))
	return (*[_MAX_ALLOC_SIZE]byte)(
		unsafe.Pointer(&buf[e.pos]))[:e.ksize]
}

func (e *leafElem) value() []byte {
	buf := (*[_MAX_ALLOC_SIZE]byte)(unsafe.Pointer(e))
	return (*[_MAX_ALLOC_SIZE]byte)(
		unsafe.Pointer(&buf[e.pos]))[e.ksize : e.ksize+e.vsize]
}

func (e *leafElem) toINode() (n *inode) {
	return &inode{
		pgid:  0,
		key:   e.key(),
		value: e.value(),
	}
}

func (e *branchElem) key() []byte {
	buf := (*[_MAX_ALLOC_SIZE]byte)(unsafe.Pointer(e))
	return (*[_MAX_ALLOC_SIZE]byte)(
		unsafe.Pointer(&buf[e.pos]))[:e.ksize:e.ksize]
}

func (e *branchElem) toINode() (n *inode) {
	return &inode{
		pgid:  e.pgid,
		key:   e.key(),
		value: nil,
	}
}
