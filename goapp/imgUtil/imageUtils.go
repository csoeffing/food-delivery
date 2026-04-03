package imgUtil

import (
	"crunchgarage/restaurant-food-delivery/util"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"slices"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/nfnt/resize"
)

func PostProcessFile(savePath string) (string, error) {
	maxWid := 256
	maxHgt := 256

	imageSize, err := GetImageSize(savePath)

	imageSizeDesc := imageSize.description()

	if err != nil {
		return "", err
	}

	if imageSize.Width > maxWid || imageSize.Height > maxHgt {
		err := ResizeImage(savePath, imageSize, ImageSize{Width: maxWid, Height: maxHgt})

		if err != nil {
			return "", err
		}

		confirmedImageSize, err := GetImageSize(savePath)

		if err == nil {
			imageSizeDesc = confirmedImageSize.description()
		}
	}

	return imageSizeDesc, nil
}

type ImageSize struct {
	Width  int
	Height int
}

func (s *ImageSize) reduce(maxSize ImageSize) ImageSize {
	xRatio := float64(maxSize.Width) / float64(s.Width)
	yRatio := float64(maxSize.Height) / float64(s.Height)

	minRatio := min(xRatio, yRatio)

	if minRatio >= 1.0 {
		return *s
	}

	newWidth := int(float64(s.Width) * minRatio)
	newHeight := int(float64(s.Height) * minRatio)

	return ImageSize{Width: newWidth, Height: newHeight}
}

func (s *ImageSize) description() string {
	return fmt.Sprintf("%dx%d", s.Width, s.Height)
}

func GetImageSize(path string) (ImageSize, error) {
	reader, err := os.Open(path)

	if err != nil {
		return ImageSize{}, err
	}

	defer reader.Close()

	im, _, err := image.DecodeConfig(reader)

	return ImageSize{Width: im.Width, Height: im.Height}, err
}

func ResizeImage(path string, imgSize, maxSize ImageSize) error {
	_, ext := util.GetFilenameAndExtension(path)

	reader, err := os.Open(path)

	if err != nil {
		return err
	}

	defer reader.Close()

	img, _, err := image.Decode(reader)

	if err != nil {
		return err
	}

	newSize := imgSize.reduce(maxSize)
	newImage := resize.Resize(uint(newSize.Width), uint(newSize.Height), img, resize.Lanczos3)

	f, err := os.Create(path)

	if err != nil {
		return err
	}

	switch ext {
	case "jpg":
		fallthrough
	case "jpeg":
		err = jpeg.Encode(f, newImage, nil)
	case "png":
		err = png.Encode(f, newImage)
	default:
		return fmt.Errorf("unknown image extension: %s", ext)
	}

	if err == nil {
		msg := fmt.Sprintf("Resized %s from (%d x %d) to (%d x %d)", path, imgSize.Width, imgSize.Height, newSize.Width, newSize.Height)
		logger.Info(msg)
	}

	return err
}

var KNOWN_IMAGE_EXTENSIONS = []string{"png", "jpeg", "jpg"}

func IsImageFile(path string) bool {
	filename := util.GetLastPathComponent(path)
	_, ext := util.GetFilenameAndExtension(filename)

	return IsKnownImageExtension(ext)
}

func IsKnownImageExtension(ext string) bool {
	return slices.Contains(KNOWN_IMAGE_EXTENSIONS, ext)
}
