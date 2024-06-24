package storage

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/google/uuid"
	"github.com/h2non/filetype"
)

func SaveImage(image []byte) (string, error) {
	uuid := uuid.New().String()

	// image가 빈 파일인가?
	if len(image) == 0 {
		return "", errors.New("empty body")
	}

	// image가 image인가?
	if !filetype.IsImage(image) {
		return "", errors.New("not image file")
	}

	// 파일 유형 검사 및 파일 확장자 추출
	kind, err := filetype.Match(image)
	if err != nil || kind == filetype.Unknown {
		return "", errors.New("unsupported file type")
	}

	err = DB().SaveExtendsion(uuid, kind.Extension)
	if err != nil {
		return "", errors.New("not save extensions")
	}

	path := fmt.Sprintf("storage/images/%s.%s", uuid, kind.Extension)
	err = os.WriteFile(path, image, 0644)
	if err != nil {
		return "", errors.New("unable to save file")
	}

	return uuid, nil
}

func FindImage(uuid string) (image.Image, error) {
	ext := DB().FindExtension(uuid)
	path := fmt.Sprintf("storage/images/%s.%s", uuid, ext)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}
