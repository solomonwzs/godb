package lmdb

import (
	"os"
	"unsafe"
)

const (
	_VERSION = 1

	_PFLAG_BRANCH   = 0x01
	_PFLAG_LEAF     = 0x02
	_PFLAG_META     = 0x04
	_PFLAG_FREELIST = 0x08

	_PGID_META_0   = 0
	_PGID_META_1   = 1
	_PGID_FREELIST = 2
	_PGID_LEAF     = 3

	_INITIAL_MMAP_SIZE = 0

	_SIZE_32K = 32 * 1024
	_SIZE_1G  = 1024 * 1024 * 1024
	_SIZE_1T  = 1024 * _SIZE_1G

	_MAX_MAP_SIZE = _SIZE_1T

	_META_CONTENT_SIZE = unsafe.Offsetof(meta{}.crc32)
)

var (
	_PAGE_SIZE int = os.Getpagesize()
)
