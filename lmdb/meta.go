package lmdb

import (
	"hash/crc32"
	"unsafe"
)

const (
	_META_CONTENT_SIZE = unsafe.Offsetof(meta{}.crc32)
)

type meta struct {
	version    uint32
	pageSize   uint32
	freelistId uint32

	crc32 uint32
}

func (m *meta) checksum() uint32 {
	crc32q := crc32.MakeTable(crc32.IEEE)
	b := (*[_META_CONTENT_SIZE]byte)(unsafe.Pointer(m))
	return crc32.Checksum(b[:], crc32q)
}

func (m *meta) init(pageSize uint32) {
	m.version = _VERSION
	m.pageSize = pageSize
	m.freelistId = _PGID_FREELIST

	m.crc32 = m.checksum()
}

func (m *meta) validate() error {
	if m.version != _VERSION {
		return ErrVersionMismatch
	} else if m.crc32 != m.checksum() {
		return ErrChecksum
	}
	return nil
}
