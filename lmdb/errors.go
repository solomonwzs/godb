package lmdb

import "errors"

var (
	ErrTimeout          error = errors.New("timeout")
	ErrVersionMismatch  error = errors.New("version mismatch")
	ErrChecksum         error = errors.New("checksum")
	ErrFileSizeTooSmall error = errors.New("file size too small")
)
