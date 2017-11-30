package lmdb

import (
	"os"
	"unsafe"
)

const (
	_VERSION = 1

	_INITIAL_MMAP_SIZE = 0

	_SIZE_32K = 32 * 1024
	_SIZE_1G  = 1024 * 1024 * 1024
	_SIZE_1T  = 1024 * _SIZE_1G

	_MAX_MAP_SIZE   = 256 * _SIZE_1T
	_MAX_ALLOC_SIZE = 2 * _SIZE_1G

	_META_CONTENT_SIZE = unsafe.Offsetof(meta{}.crc32)

	_MAX_ELEMENT_COUNT = 0xffffffff
)

var (
	_PAGE_SIZE int = os.Getpagesize()
)
