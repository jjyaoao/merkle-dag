package merkledag

type TestFile struct {
	name string
	data []byte
}

func (file *TestFile) Size() uint64 {
	return uint64(len(file.data))
}

func (file *TestFile) Name() string {
	return file.name
}

func (file *TestFile) Type() int {
	return FILE
}

func (file *TestFile) Bytes() []byte {
	return file.data
}

type testDirIter struct {
	list []Node
	iter int
}

func (iter *testDirIter) Next() bool {
	if iter.iter+1 < len(iter.list) {
		iter.iter += 1
		return true
	}
	return false
}

func (iter *testDirIter) Node() Node {
	return iter.list[iter.iter]
}

type TestDir struct {
	list []Node
	name string
}

func (dir *TestDir) Size() uint64 {
	var len uint64 = 0
	for i := range dir.list {
		len += dir.list[i].Size()
	}
	return len
}

func (dir *TestDir) Name() string {
	return dir.name
}

func (dir *TestDir) Type() int {
	return DIR
}

func (dir *TestDir) It() DirIterator {
	it := &testDirIter{
		list: dir.list,
		iter: -1,
	}
	return it
}

type HashMap struct {
	mp map[string]([]byte)
}

func (hmp *HashMap) Has(key []byte) (bool, error) {
	return hmp.mp[string(key)] != nil, nil
}

func (hmp *HashMap) Put(key, value []byte) error {
	flag, _ := hmp.Has(key)
	if flag {
		panic("Key is same")
	}
	hmp.mp[string(key)] = value
	return nil
}

func (hmp *HashMap) Get(key []byte) ([]byte, error) {
	flag, _ := hmp.Has(key)
	if !flag {
		panic("Don't have the key")
	}
	return hmp.mp[string(key)], nil
}

func (hmp *HashMap) Delete(key []byte) error {
	return nil
}
