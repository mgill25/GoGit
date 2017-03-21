package main

import (
	"crypto/sha1"
	"fmt"
)

type Blob struct {
	hash []byte
}

// Pattern from https://play.golang.org/p/YUaWWEeB4U
func (b *Blob) createBlob(content string) {
	hash := sha1.New()
	hash.Write([]byte(content))
	byteSlice := hash.Sum(nil)
	b.hash = byteSlice
	fmt.Printf("%x\n", b.hash)
}
