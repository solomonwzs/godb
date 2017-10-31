package lmdb

import (
	"syscall"
	"time"
	"unsafe"
)

func flock(db *DB, timeout time.Duration) (err error) {
	start := time.Now()

	for {
		if timeout > 0 && time.Since(start) > timeout {
			return ErrTimeout
		}

		flag := syscall.LOCK_SH
		if db.readOnly {
			flag = syscall.LOCK_EX
		}

		err = syscall.Flock(int(db.file.Fd()), flag|syscall.LOCK_NB)
		if err == nil || err != syscall.EWOULDBLOCK {
			return
		}

		time.Sleep(50 * time.Millisecond)
	}

	return ErrTimeout
}

func mmap(db *DB, size int) (err error) {
	var b []byte
	if b, err = syscall.Mmap(int(db.file.Fd()), 0, size,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED); err != nil {
		return
	}

	db.dataref = b
	db.data = (*[_MAX_MAP_SIZE]byte)(unsafe.Pointer(&b[0]))
	db.dataSize = size

	return
}

func munmap(db *DB) (err error) {
	err = nil
	if db.dataref == nil {
		return
	}

	if err = syscall.Munmap(db.dataref); err != nil {
		return
	}

	db.dataref = nil
	db.data = nil
	db.dataSize = 0

	return
}
