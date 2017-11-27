package lmdb

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

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

	path := db.Path()
	fd, _ := os.OpenFile(path, os.O_RDONLY, 0666)
	defer func() {
		fd.Close()
		// os.Remove(path)
	}()

	bucketName := []byte("bucket")
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			if b, err = tx.CreateBucket(bucketName); err != nil {
				t.Fatal(err)
			}
		}

		for i := 0; i < 64; i++ {
			key := []byte(fmt.Sprintf("\x02\x04\x08-%d", i))
			valueSize := 3 * 1024
			value := make([]byte, valueSize)
			b.Put(key, value)
		}

		return nil
	}); err != nil {
		t.Fatal(err)
	}
	info, _ := fd.Stat()
	fmt.Printf("%+v\n", info.Size())
}

func _TestMmap(t *testing.T) {
	map_file, err := os.Create("/tmp/test3")
	// map_file, err := os.OpenFile("/tmp/test2", os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	info, _ := map_file.Stat()
	sz := int(info.Size())

	mmap, err := syscall.Mmap(int(map_file.Fd()), 0, sz,
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// _, err = map_file.Seek(int64(s-1), 0)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// _, err = map_file.Write([]byte(" "))
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// map_array := (*[sz]int)(unsafe.Pointer(&mmap[0]))
	page := 4
	pageSize := 4096
	fmt.Println(mmap[page*pageSize : (page+1)*pageSize])

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
