package lmdb

import (
	"os"
	"unsafe"
)

const (
	_VERSION = 1

	_PAGE_FLAG_BRANCH   = 0x01
	_PAGE_FLAG_LEAF     = 0x02
	_PAGE_FLAG_META     = 0x04
	_PAGE_FLAG_FREELIST = 0x08

	_INITIAL_MMAP_SIZE = 0

	_SIZE_32K = 32 * 1024
	_SIZE_1G  = 1024 * 1024 * 1024
	_SIZE_1T  = 1024 * _SIZE_1G

	_MAX_MAP_SIZE   = 256 * _SIZE_1T
	_MAX_ALLOC_SIZE = 2 * _SIZE_1G

	_META_CONTENT_SIZE = unsafe.Offsetof(meta{}.crc32)
)

var (
	_PAGE_SIZE int = os.Getpagesize()
)
