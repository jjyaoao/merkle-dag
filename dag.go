package merkledag

import (
	"encoding/json"
	"hash"
)

type Link struct {
	Name string
	Hash []byte
	Size int
}

type Object struct {
	Links []Link
	Data  []byte
}

func dfsForSliceFile(hight int, node File, store KVStore, seedId int, h hash.Hash) (*Object, int) {

	if hight == 1 {
		if (len(node.Bytes()) - seedId) <= 256*1024 {
			data := node.Bytes()[seedId:]
			blob := Object{
				Links: nil,
				Data:  data,
			}
			jsonMarshal, _ := json.Marshal(blob)
			h.Write(jsonMarshal)
			store.Put(h.Sum(nil), data)
			return &blob, len(data)
		}
		links := &Object{}
		lenData := 0
		for i := 1; i <= 4096; i++ {
			end := seedId + 256*1024
			if len(node.Bytes()) < end {
				end = len(node.Bytes())
			}
			data := node.Bytes()[seedId:end]
			blob := Object{
				Links: nil,
				Data:  data,
			}
			lenData += len(data)
			jsonMarshal, _ := json.Marshal(blob)
			h.Write(jsonMarshal)
			store.Put(h.Sum(nil), data)
			links.Links = append(links.Links, Link{
				Hash: h.Sum(nil),
				Size: len(data),
			})
			links.Data = append(links.Data, []byte("blob")...)
			seedId += 256 * 1024
			if seedId >= len(node.Bytes()) {
				break
			}
		}
		jsonMarshal, _ := json.Marshal(links)
		h.Write(jsonMarshal)
		store.Put(h.Sum(nil), jsonMarshal)
		return links, lenData
	} else {
		links := &Object{}
		lenData := 0
		for i := 1; i <= 4096; i++ {
			if seedId >= len(node.Bytes()) {
				break
			}
			tmp, lens := dfsForSliceFile(hight-1, node, store, seedId, h)
			lenData += lens
			jsonMarshal, _ := json.Marshal(tmp)
			h.Write(jsonMarshal)
			links.Links = append(links.Links, Link{
				Hash: h.Sum(nil),
				// 对于link 它size是什么？
				Size: lens,
			})
			typeName := "link"
			if tmp.Links == nil {
				typeName = "blob"
			}
			links.Data = append(links.Data, []byte(typeName)...)
		}
		jsonMarshal, _ := json.Marshal(links)
		h.Write(jsonMarshal)
		store.Put(h.Sum(nil), jsonMarshal)
		return links, lenData
	}
}
func sliceFile(node File, store KVStore, h hash.Hash) *Object {
	// fmt.Println("222222")
	if len(node.Bytes()) <= 256*1024 {
		data := node.Bytes()
		blob := Object{
			Links: nil,
			Data:  data,
		}
		jsonMarshal, _ := json.Marshal(blob)
		h.Write(jsonMarshal)
		store.Put(h.Sum(nil), data)
		return &blob
	}
	linkLen := (len(node.Bytes()) + (256*1024 - 1)) / (256 * 1024)
	hight := 0
	tmp := linkLen
	for {
		hight++
		tmp /= 4096
		// fmt.Println(tmp)
		if tmp == 0 {
			break
		}
	}
	// fmt.Println(hight)
	res, _ := dfsForSliceFile(hight, node, store, 0, h)
	return res
}

func sliceDir(node Dir, store KVStore, h hash.Hash) *Object {
	iter := node.It()
	treeObject := &Object{}
	for iter.Next() {
		node := iter.Node()
		if node.Type() == FILE {
			file := node.(File)
			tmp := sliceFile(file, store, h)
			jsonMarshal, _ := json.Marshal(tmp)
			h.Write(jsonMarshal)
			treeObject.Links = append(treeObject.Links, Link{
				Hash: h.Sum(nil),
				Size: int(file.Size()),
				Name: file.Name(),
			})
			typeName := "link"
			if tmp.Links == nil {
				typeName = "blob"
			}
			treeObject.Data = append(treeObject.Data, []byte(typeName)...)
		} else {
			dir := node.(Dir)
			tmp := sliceDir(dir, store, h)
			jsonMarshal, _ := json.Marshal(tmp)
			h.Write(jsonMarshal)
			treeObject.Links = append(treeObject.Links, Link{
				Hash: h.Sum(nil),
				Size: int(dir.Size()),
				Name: dir.Name(),
			})
			typeName := "tree"
			treeObject.Data = append(treeObject.Data, []byte(typeName)...)
		}
	}
	jsonMarshal, _ := json.Marshal(treeObject)
	h.Write(jsonMarshal)
	store.Put(h.Sum(nil), jsonMarshal)
	return treeObject
}
func Add(store KVStore, node Node, h hash.Hash) []byte {
	// TODO 将分片写入到KVStore中，并返回Merkle Root
	if node.Type() == FILE {
		file := node.(File)
		tmp := sliceFile(file, store, h)
		jsonMarshal, _ := json.Marshal(tmp)
		h.Write(jsonMarshal)
		return h.Sum(nil)
	} else {
		dir := node.(Dir)
		tmp := sliceDir(dir, store, h)
		jsonMarshal, _ := json.Marshal(tmp)
		h.Write(jsonMarshal)
		return h.Sum(nil)
	}
}