package lmdb

import (
	"os"
	"sync"
	"time"
	"unsafe"
)

type DB struct {
	path     string
	file     *os.File
	opened   bool
	readOnly bool
	pageSize int

	dataref  []byte
	data     *[_MAX_MAP_SIZE]byte
	dataSize int

	meta0 *meta
	meta1 *meta

	metaLock sync.Mutex
	mmapLock sync.RWMutex
}

type Options struct {
	ReadOnly bool
	FileMode os.FileMode
	TimeOut  time.Duration
}

func Open(path string, opt *Options) (db *DB, err error) {
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

	var info os.FileInfo
	if info, err = db.file.Stat(); err != nil {
		return
	} else if info.Size() == 0 {
		if err = db.init(); err != nil {
			return
		}
	} else {
		buf := make([]byte, _PAGE_SIZE*2)
		if _, err = db.file.ReadAt(buf, 0); err != nil {
			return
		}

		m := db.getPageFromBytes(buf, 0).getMeta()
		if err = m.validate(); err != nil {
			db.close()
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
	buf := make([]byte, db.pageSize*4)

	for i := pageid(0); i < 2; i++ {
		p := db.getPageFromBytes(buf, i)
		p.pgid = i
		p.flags = _PAGE_FLAG_META

		m := p.getMeta()
		m.version = _VERSION
		m.pageSize = uint32(db.pageSize)
		m.freelistId = 2
		m.tid = txid(i)
		m.root = 3

		m.crc32 = m.checksum()
	}

	p := db.getPageFromBytes(buf, 2)
	p.pgid = 2
	p.count = 0
	p.overflow = 0
	p.flags = _PAGE_FLAG_FREELIST

	p = db.getPageFromBytes(buf, 3)
	p.pgid = 3
	p.count = 0
	p.overflow = 0
	p.flags = _PAGE_FLAG_LEAF

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

func (db *DB) getPageFromBytes(b []byte, id pageid) *page {
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

	db.meta0 = db.getPage(0).getMeta()
	db.meta1 = db.getPage(1).getMeta()

	if err = db.meta0.validate(); err != nil {
		return
	}
	if err = db.meta1.validate(); err != nil {
		return
	}

	return
}

func (db *DB) getPage(pgid pageid) *page {
	return (*page)(unsafe.Pointer(&db.data[int(pgid)*db.pageSize]))
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
