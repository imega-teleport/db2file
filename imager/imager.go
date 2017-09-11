package imager

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
)

type ImageInfo struct {
	Width  int
	Height int
	Mime   string
	Name   string
}

func GetImageInfo(filename string) ImageInfo {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	defer file.Close()

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}

	return ImageInfo{
		Width:  image.Width,
		Height: image.Height,
	}
}
