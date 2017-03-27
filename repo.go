package main

import (
	"path/filepath"
	"os"
)

type Repo struct {
	mainPath string
	objPath string
}

func (r *Repo) createObjectDatabase() {
	dirName, err := filepath.Abs(r.mainPath + "objects")
	check(err)
	os.Mkdir(dirName, 0755)
	r.objPath = dirName
}
