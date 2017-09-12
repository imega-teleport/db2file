package imager

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/h2non/filetype"
)

type ImageInfo struct {
	Width  int
	Height int
	Mime   string
	Name   string
}

func GetImageInfo(filename string) (ImageInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return ImageInfo{}, err
	}
	defer file.Close()

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		return ImageInfo{}, err
	}

	file.Seek(0, 0)
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return ImageInfo{}, err
	}

	kind, err := filetype.Match(buffer[:n])
	if err != nil {
		return ImageInfo{}, err
	}

	return ImageInfo{
		Width:  image.Width,
		Height: image.Height,
		Mime:   kind.MIME.Value,
		Name:   file.Name(),
	}, nil
}
