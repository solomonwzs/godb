package lmdb

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

func NewNode(p *page, parent *node) (n *node, err error) {
	if p.flags != _PAGE_FLAG_BRANCH || p.flags != _PAGE_FLAG_LEAF {
		return nil, ErrPageFlags
	}

	n = &node{}
	if p.flags == _PAGE_FLAG_BRANCH {
		elemList := p.getBranchElems()
		n.ilist = make([]*inode, len(elemList))

		for i, elem := range elemList {
			n.ilist[i] = new(inode)
			n.ilist[i].value = nil
			n.pgid = elem.pgid

			if elem.ksize != 0 {
				n.ilist[i].key = elem.key()
			}
		}
	}

	return
}
