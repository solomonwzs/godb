package lmdb

import (
	"unsafe"
)

type node struct {
	pgid   pageid
	parent *node
	ilist  []*inode
	isLeaf bool
}

type inode struct {
	pgid  pageid
	key   []byte
	value []byte
}

func (n *node) pageElemSize() int {
	if n.isLeaf {
		return len(n.ilist) * _LEAF_ELEM_SIZE
	} else {
		return len(n.ilist) * _BRANCH_ELEM_SIZE
	}
}

func (n *node) pageContentSize() int {
	sz := 0
	if n.isLeaf {
		for _, i := range n.ilist {
			sz += len(i.key) + len(i.value)
		}
	} else {
		for _, i := range n.ilist {
			sz += len(i.key)
		}
		sz += len(n.ilist) * int(unsafe.Sizeof(inode{}.pgid))
	}

	return sz
}

func (n *node) readFrom(p *page) (err error) {
	if p.flags&_PAGE_FLAG_BRANCH == 0 && p.flags&_PAGE_FLAG_LEAF == 0 {
		return ErrPageFlags
	}

	n.pgid = p.pgid
	if p.flags&_PAGE_FLAG_BRANCH != 0 {
		el := p.getBranchElems()
		n.isLeaf = false
		n.ilist = make([]*inode, len(el))

		for i := 0; i < len(el); i++ {
			n.ilist[i] = el[i].toINode()
		}
	} else if p.flags&_PAGE_FLAG_LEAF != 0 {
		el := p.getLeafElems()
		n.isLeaf = true
		n.ilist = make([]*inode, len(el))

		for i := 0; i < len(el); i++ {
			n.ilist[i] = el[i].toINode()
		}
	}

	return
}

func (n *node) writeTo(p *page) (err error) {
	if n.isLeaf {
		p.flags |= _PAGE_FLAG_LEAF
	} else {
		p.flags |= _PAGE_FLAG_BRANCH
	}

	count := len(n.ilist)
	if count > _MAX_NODE_COUNT {
		return ErrSizeOverflow
	}
	p.count = uint16(count)

	if p.count == 0 {
		return
	}

	esz := n.pageElemSize()
	csz := n.pageContentSize()
	buf := (*[_MAX_ALLOC_SIZE]byte)(
		unsafe.Pointer(&p.ptr))[esz : esz+csz]
	if n.isLeaf {
		el := p.getLeafElems()
		for i, ni := range n.ilist {
			el[i].ksize = uint32(len(ni.key))
			el[i].vsize = uint32(len(ni.value))
			el[i].pos = uint32(uintptr(unsafe.Pointer(&buf[0])) -
				uintptr(unsafe.Pointer(&el[i])))

			copy(buf, ni.key)
			copy(buf[el[i].ksize:], ni.value)
			buf = buf[el[i].ksize+el[i].vsize:]
		}
	} else {
		el := p.getBranchElems()
		for i, ni := range n.ilist {
			el[i].ksize = uint32(len(ni.key))
			el[i].pgid = ni.pgid
			el[i].pos = uint32(uintptr(unsafe.Pointer(&buf[0])) -
				uintptr(unsafe.Pointer(&el[i])))

			copy(buf, ni.key)
			buf = buf[el[i].ksize:]
		}
	}

	return
}
