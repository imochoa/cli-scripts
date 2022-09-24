package main

import (
	"log"

	// "github.com/atotto/clipboard"
	"golang.design/x/clipboard"
)

// https://pkg.go.dev/github.com/atotto/clipboard#section-readme

// Handle images?
// https://pkg.go.dev/golang.design/x/clipboard
// https://stackoverflow.com/questions/6841532/linux-image-from-clipboard

func main() {

	// str, _ := clipboard.ReadAll()

	// Init returns an error if the package is not ready for use.
	err := clipboard.Init()
	if err != nil {
		log.Printf("PANIC!")
		panic(err)
	}
	//The most common operations are Read and Write. To use them:

	// write/read text format data of the clipboard, and
	// the byte buffer regarding the text are UTF8 encoded.
	// clipboard.Write(clipboard.FmtText, []byte("text data"))
	str := clipboard.Read(clipboard.FmtText)
	log.Printf("clipboard STR: %s", str)

	// write/read image format data of the clipboard, and
	// the byte buffer regarding the image are PNG encoded.
	// clipboard.Write(clipboard.FmtImage, []byte("image data"))
	imgData := clipboard.Read(clipboard.FmtImage)

	log.Printf("clipboard: %s", imgData)

}
