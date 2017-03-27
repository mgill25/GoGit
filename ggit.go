package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
	"fmt"
)

func main() {
	app := cli.NewApp()
	app.Name = "ggit"
	app.Usage = "Git - Implemented in GoLang!"
	app.Commands = []cli.Command{
		// Plumbing Commands
		{
			Name:  "hash-object",
			Usage: "Compute hash ID and optionally create a blob from the file",
			Action: func(c *cli.Context) {
				hashObject(c.Args().First())
			},
		},
		{
			Name: "cat-file",
			Usage: "Retrieve the contents of an object",
			Action: func(c *cli.Context) {
				catFile(c.Args().First())
			},
		},
		// Porcelain Commands
		{
			Name: "init",
			Usage: "Initialize a new Git repository",
			Action: func(c *cli.Context) {
				initRepo()
			},
		},
	}
	app.Run(os.Args)
}

func hashObject(fileName string) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	filePath := dir + "/" + fileName
	fileData := readContentFile(filePath)
	dataWithHeaders := prependContentHeaders("blob", fileData)

	// Now that we have all the data with headers, we pass it through sha1 to get blob id
	b := Blob{}
	b.createBlob(dataWithHeaders)
}

func catFile(fileName string) {
	repo := getCurrentRepo()
	// Split up the file name:
	dirName := fileName[:2]
	objFileName := fileName[2:]
	dir, err := filepath.Abs(repo.objPath + "/" + dirName)
	if err != nil {
		log.Fatal(err)
	}
	objectHash := dir + objFileName
	filePath := dir + "/" + objFileName

	// Initialize the blob object
	b := Blob{hash: []byte(objectHash), path: filePath}
	b.readBlob()
}

// Initialize a new git repository
func initRepo() {
	dirName, err := filepath.Abs(filepath.Dir(os.Args[0]) + "git")
	check(err)
	os.Mkdir(dirName, 0755)
	fmt.Println("Initialized empty Git repository in ", dirName)
}

func getCurrentRepo() Repo {
	dirName, err := filepath.Abs(filepath.Dir(os.Args[0]) + "git")
	check(err)
	objPath, err := filepath.Abs(dirName + "/objects")
	if _, err := os.Stat(objPath); os.IsNotExist(err) {
		os.Mkdir(objPath, 0755)
	}
	return Repo{mainPath: dirName, objPath: objPath}
}
