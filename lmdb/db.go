package lmdb

import (
	"os"
	"sync"
	"time"
	"unsafe"
)

const (
	_VERSION = 1

	_INITIAL_MMAP_SIZE = 0

	_SIZE_32K = 32 * 1024
	_SIZE_1G  = 1024 * 1024 * 1024
	_SIZE_1T  = 1024 * _SIZE_1G

	_MAX_MMAP_SIZE = _SIZE_1T
)

var (
	_PAGE_SIZE int = os.Getpagesize()
)

type DB struct {
	path     string
	file     *os.File
	opened   bool
	readOnly bool
	pageSize int

	data     *[]byte
	dataref  []byte
	dataSize int

	meta *meta

	mmapLock sync.RWMutex
}

type Options struct {
	ReadOnly bool
	FileMode os.FileMode
	TimeOut  time.Duration
}

func Open(path string, opt *Options) (db *DB, err error) {
	var (
		info os.FileInfo
	)

	db = &DB{
		path:   path,
		file:   nil,
		opened: true,
	}

	flag := os.O_RDWR
	if opt.ReadOnly {
		flag = os.O_RDONLY
		db.readOnly = true
	}
	if db.file, err = os.OpenFile(path, flag|os.O_CREATE,
		opt.FileMode); err != nil {
		db.close()
		return
	}

	if err = flock(db, opt.TimeOut); err != nil {
		db.close()
		return
	}

	if info, err = db.file.Stat(); err != nil {
		return
	} else if info.Size() == 0 {
		if err = db.init(); err != nil {
			return
		}
	} else {
		buf := make([]byte, _PAGE_SIZE)
		if _, err = db.file.ReadAt(buf, 0); err != nil {
			return
		}

		m := db.getPageFromBuffer(buf, _PGID_META).getMeta()
		if err = m.validate(); err != nil {
			return
		}
		db.pageSize = int(m.pageSize)
	}

	if err = db.mmap(); err != nil {
		return
	}

	return
}

func (db *DB) init() (err error) {
	db.pageSize = os.Getpagesize()
	buf := make([]byte, db.pageSize*3)

	p := db.getPageFromBuffer(buf, _PGID_META)
	p.id = _PGID_META
	p.flags = _PFLAG_META
	m := p.getMeta()
	m.init(uint32(db.pageSize))

	p = db.getPageFromBuffer(buf, _PGID_FREELIST)
	p.id = _PGID_FREELIST
	p.flags = _PFLAG_FREELIST

	p = db.getPageFromBuffer(buf, _PGID_LEAF)
	p.id = _PGID_LEAF
	p.flags = _PFLAG_LEAF

	if _, err = db.file.Write(buf); err != nil {
		return
	}
	if err = db.file.Sync(); err != nil {
		return
	}

	return
}

func (db *DB) close() (err error) {
	err = nil
	if !db.opened {
		return
	}

	if err = munmap(db); err != nil {
		return
	}

	if db.file != nil {
		if err = db.file.Close(); err != nil {
			return
		}
	}

	db.file = nil
	db.path = ""
	db.opened = false

	return
}

func (db *DB) getPageFromBuffer(b []byte, id pageid) *page {
	return (*page)(unsafe.Pointer(&b[int(id)*db.pageSize]))
}

func (db *DB) mmap() (err error) {
	db.mmapLock.Lock()
	defer db.mmapLock.Unlock()

	var (
		info os.FileInfo
		size int
	)

	info, err = db.file.Stat()
	if err != nil {
		return
	} else if int(info.Size()) < db.pageSize {
		return ErrFileSizeTooSmall
	}

	size = determineMmapSize(int(info.Size()), db.pageSize)
	if err = mmap(db, size); err != nil {
		return err
	}

	m := db.getPage(_PGID_META).getMeta()
	if err = m.validate(); err != nil {
		return
	}

	return
}

func (db *DB) getPage(pgid pageid) *page {
	return (*page)(unsafe.Pointer(&db.dataref[int(pgid)*db.pageSize]))
}

func determineMmapSize(size int, pageSize int) (newSize int) {
	newSize = _SIZE_32K
	for newSize < _SIZE_1G {
		if size <= newSize {
			return
		}
		newSize <<= 1
	}

	if newSize%pageSize != 0 {
		newSize = (newSize/pageSize + 1) * pageSize
	}

	return
}
