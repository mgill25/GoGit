package main

import (
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
)

type Blob struct {
	hash []byte
	path string
}

func (b *Blob) createBlob(content string) {
	// First hash the content
	b.hashContent(content)
	// Then create a compressed blob object
	dirName := hex.EncodeToString(b.hash[:1])
	objName := hex.EncodeToString(b.hash[1:])
	os.Mkdir(dirName, 0755)
	objFile, err := os.Create(dirName + string("/") + objName)
	check(err)
	w := zlib.NewWriter(objFile)
	w.Write([]byte(content))
	w.Close()
	//fmt.Println(b.Bytes())
}

// Pattern from https://play.golang.org/p/YUaWWEeB4U
func (b *Blob) hashContent(content string) {
	hash := sha1.New()
	hash.Write([]byte(content))
	byteSlice := hash.Sum(nil)
	b.hash = byteSlice
	fmt.Printf("%x\n", b.hash)
}
