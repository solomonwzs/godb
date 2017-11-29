package lmdb

import "errors"

var (
	ErrTimeout          = errors.New("lmdb: timeout")
	ErrVersionMismatch  = errors.New("lmdb: version mismatch")
	ErrChecksum         = errors.New("lmdb: checksum")
	ErrFileSizeTooSmall = errors.New("lmdb: file size too small")
	ErrBytesLen         = errors.New("lmdb: bytes length")
	ErrPageFlags        = errors.New("lmdb: page flags error")
)
