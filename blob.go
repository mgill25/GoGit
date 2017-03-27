package main

import (
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"io"
	"bytes"
	"strings"
)

type Blob struct {
	hash []byte
	path string
}

func (b *Blob) createBlob(content string) {
	// First hash the content
	b.hashContent(content)
	// Then create a compressed blob object
	prefixName:= hex.EncodeToString(b.hash[:1])
	objName := hex.EncodeToString(b.hash[1:])
	r := getCurrentRepo()
	os.Mkdir(r.objPath + string("/") + prefixName, 0755)
	objFile, err := os.Create(r.objPath + string("/") + prefixName + string("/") + objName)
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

func (b *Blob) readBlob() {
	objFile, err := os.Open(b.path)
	check(err)
	var out bytes.Buffer
	r, _ := zlib.NewReader(objFile)
	io.Copy(&out, r)
	data := string(out.Bytes())
	split_data := strings.Split(data, "\000")
	//headers := split_data[0]
	content := split_data[1]
	os.Stdout.Write([]byte(content))
}
