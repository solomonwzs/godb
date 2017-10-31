package lmdb

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"
	"unsafe"

	"github.com/boltdb/bolt"
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

func TestBolt(t *testing.T) {
	db, err := bolt.Open("/tmp/test2", 0666, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("widgets"))
		if err != nil {
			t.Fatal(err)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestMmap(t *testing.T) {
	const n = 1e3
	s := int(unsafe.Sizeof(0)) * n

	map_file, err := os.Create("/tmp/test3")
	// map_file, err := os.OpenFile("/tmp/test1", os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mmap, err := syscall.Mmap(int(map_file.Fd()), 0, 10,
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = map_file.Seek(int64(s-1), 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = map_file.Write([]byte(" "))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	map_array := (*[n]int)(unsafe.Pointer(&mmap[0]))
	fmt.Println(map_array[:100])

	// for i := 0; i < n; i++ {
	// 	map_array[i] = i * i
	// }

	// err = syscall.Munmap(mmap)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// err = map_file.Close()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
}
