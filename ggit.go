package main

import (
	"os"
	"github.com/urfave/cli"
	"path/filepath"
	"log"
)


func main() {
	app := cli.NewApp()
	app.Name = "ggit"
	app.Usage = "Git - Implemented in GoLang!"
	app.Commands = []cli.Command{
		{
			Name: "hash-object",
			Usage: "Compute hash ID and optionally create a blob from the file",
			Action: func(c *cli.Context) {
				hashObject(c.Args().First())
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
