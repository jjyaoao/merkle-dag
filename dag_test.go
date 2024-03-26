package merkledag

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

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

func TestDagToFile(t *testing.T) {
	hashMapInstance := &HashMap{
		mp: make(map[string][]byte),
	}
	hashFunc := sha256.New()

	// 定义文件夹
	folderPath := "D:\\Information\\作业=-=\\分布式\\merkle-dag"
	filesInDir, _ := ioutil.ReadDir(folderPath)

	directory := &TestDir{
		list: make([]Node, len(filesInDir)),
		name: "/",
	}

	for idx, fileItem := range filesInDir {
		filePath := folderPath + "/" + fileItem.Name()

		if fileItem.IsDir() {
			dirContent := explore(filePath)
			dirContent.name = fileItem.Name()
			directory.list[idx] = dirContent
		} else {
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatal(err)
			}
			fileNode := &TestFile{
				name: fileItem.Name(),
				data: fileContent,
			}
			directory.list[idx] = fileNode
		}
	}

	rootHash := Add(hashMapInstance, directory, hashFunc)
	fmt.Printf("%x\n", rootHash)

	// Retrieve file from the DAG
	bufferGo := Hash2File(hashMapInstance, rootHash, "/pkg/mod/bazil.org/fuse@v0.0.0-20200117225306-7b5117fecadc/buffer.go", nil)
	fmt.Println(string(bufferGo))

	// Retrieve a new folder from the DAG
	newFolderContent := Hash2File(hashMapInstance, rootHash, "newfolder", nil)
	originalContent, _ := os.ReadFile("D:\\Information\\作业=-=\\分布式\\merkle-dag\\newfolder")

	hashFunc.Reset()
	hashFunc.Write(originalContent)
	hash1 := hashFunc.Sum(nil)

	hashFunc.Reset()
	hashFunc.Write(newFolderContent)
	hash2 := hashFunc.Sum(nil)

	fmt.Println(hash1)
	fmt.Println(hash2)
	fmt.Println(string(hash1) == string(hash2))
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
