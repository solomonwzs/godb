package lmdb

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

func TestBase(t *testing.T) {
	db, err := Open("/tmp/test1", &Options{
		FileMode: 0666,
		ReadOnly: false,
		TimeOut:  5 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(&db)
	}
}

func TestMMAP(t *testing.T) {
	path := "/tmp/test-mmap"
	mapFile, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		mapFile.Close()
		os.Remove(path)
	}()

	size := 64
	mmap, err := syscall.Mmap(int(mapFile.Fd()), 0, size,
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		t.Fatal(err)
	}

	_, err = mapFile.Seek(int64(size-1), 0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = mapFile.Write([]byte{0x00})
	if err != nil {
		t.Fatal(err)
	}

	debugln(mmap)
}

func TestNode(t *testing.T) {
	buf := make([]byte, 1024, 1024)
	p := (*page)(unsafe.Pointer(&buf[0]))

	n := &node{
		pgid:   123,
		parent: nil,
		ilist:  make([]*inode, 3, 3),
		isLeaf: true,
	}

	for i := 0; i < len(n.ilist); i++ {
		n.ilist[i] = &inode{
			pgid:  pageid(i),
			key:   []byte(fmt.Sprintf("key-%d", i)),
			value: []byte(fmt.Sprintf("value-%d", i)),
		}
	}
	n.writeTo(p)

	n1 := new(node)
	err := n1.readFrom(p)
	if err != nil {
		t.Fatal(err)
	}
	for i, _ := range n1.ilist {
		debugln(n1.ilist[i])
	}
}
