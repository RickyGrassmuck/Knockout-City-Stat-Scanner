package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"path/filepath"
)

func getFiles() ([]string, error) {
	jpgFiles, err := filepath.Glob("*.j*g")
	if err != nil {
		return nil, err
	}

	pngFiles, err := filepath.Glob("*.png")
	if err != nil {
		return nil, err
	}

	allFiles := append(jpgFiles, pngFiles...)
	return allFiles, nil
}

// Guess image format from gif/jpeg/png/webp
func guessImageFormat(r io.Reader) (format string, err error) {
	_, format, err = image.DecodeConfig(r)
	return
}

// Guess image mime types from gif/jpeg/png/webp
func guessImageMimeTypes(r io.Reader) string {
	format, _ := guessImageFormat(r)
	if format == "" {
		return ""
	}
	return mime.TypeByExtension("." + format)
}
