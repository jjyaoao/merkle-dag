package merkledag

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

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

func TestDagStructure(t *testing.T) {
	store := &HashMap{
		mp: make(map[string][]byte),
	}
	hasher := sha256.New()
	// 一个小文件的测试
	smallFile := &TestFile{
		name: "tiny",
		data: []byte("这是一个用于测试的小文件"),
	}
	rootHash := Add(store, smallFile, hasher)
	fmt.Printf("%x\n", rootHash)

	// 一个大文件的测试
	store = &HashMap{
		mp: make(map[string][]byte),
	}
	hasher.Reset()
	bigFileContent, err := os.ReadFile("D:\\Information\\作业=-=\\分布式\\merkle-dag\\213_2021131120_陈思州_1.rar")
	if err != nil {
		t.Error(err)
	}

	bigFile := &TestFile{
		name: "large",
		data: bigFileContent,
	}

	rootHash = Add(store, bigFile, hasher)
	fmt.Printf("%x\n", rootHash)

	// 一个文件夹的测试
	store = &HashMap{
		mp: make(map[string][]byte),
	}
	hasher.Reset()
	dirPath := "D:\\Information\\作业=-=\\分布式\\merkle-dag"
	entries, _ := ioutil.ReadDir(dirPath)
	directory := &TestDir{
		list: make([]Node, len(entries)),
		name: "Docs",
	}
	for i, entry := range entries {
		entryPath := dirPath + "/" + entry.Name()
		if entry.IsDir() {
			subDir := explore(entryPath)
			subDir.name = entry.Name()
			directory.list[i] = subDir
		} else {
			fileContent, err := os.ReadFile(entryPath)
			if err != nil {
				t.Fatal(err)
			}
			entryFile := &TestFile{
				name: entry.Name(),
				data: fileContent,
			}
			directory.list[i] = entryFile
		}
	}
	rootHash = Add(store, directory, hasher)
	fmt.Printf("%x\n", rootHash)
}

func explore(dirPath string) *TestDir {
	entries, _ := ioutil.ReadDir(dirPath)
	directory := &TestDir{
		list: make([]Node, len(entries)),
	}
	for i, entry := range entries {
		entryPath := dirPath + "/" + entry.Name()
		if entry.IsDir() {
			subDir := explore(entryPath)
			subDir.name = entry.Name()
			directory.list[i] = subDir
		} else {
			fileContent, err := os.ReadFile(entryPath)
			if err != nil {
				subDir := explore(entryPath)
				subDir.name = entry.Name()
				directory.list[i] = subDir
				continue
			}
			entryFile := &TestFile{
				name: entry.Name(),
				data: fileContent,
			}
			directory.list[i] = entryFile
		}
	}
	return directory
}
