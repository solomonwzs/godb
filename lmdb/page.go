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
	if p.flags != _PAGE_FLAG_LEAF || p.count == 0 {
		return nil
	}
	return (*[_MAX_ELEMENT_COUNT]leafElem)(
		unsafe.Pointer(&p.ptr))[:p.count]
}

func (p *page) getBranchElems() []branchElem {
	if p.flags != _PAGE_FLAG_BRANCH || p.count == 0 {
		return nil
	}
	return (*[_MAX_ELEMENT_COUNT]branchElem)(
		unsafe.Pointer(&p.ptr))[:p.count]
}

func (e *leafElem) key() []byte {
	buf := (*[_MAX_ALLOC_SIZE]byte)(unsafe.Pointer(e))
	return (*[_MAX_ALLOC_SIZE]byte)(
		unsafe.Pointer(&buf[e.pos]))[:e.ksize:e.ksize]
}

func (e *leafElem) value() []byte {
	buf := (*[_MAX_ALLOC_SIZE]byte)(unsafe.Pointer(e))
	return (*[_MAX_ALLOC_SIZE]byte)(
		unsafe.Pointer(&buf[e.pos]))[e.ksize : e.ksize+e.vsize : e.vsize]
}

func (e *branchElem) key() []byte {
	buf := (*[_MAX_ALLOC_SIZE]byte)(unsafe.Pointer(e))
	return (*[_MAX_ALLOC_SIZE]byte)(
		unsafe.Pointer(&buf[e.pos]))[:e.ksize:e.ksize]
}
