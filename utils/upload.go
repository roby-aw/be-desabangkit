package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func Upload(name, email string, file *multipart.FileHeader) error {
	// Read form fields

	//-----------
	// Read file
	//-----------

	// Source
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fileLocation := filepath.Join(dir, "utils/img", file.Filename)

	// Destination
	dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()
	tjn, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	// Copy
	if _, err = io.Copy(tjn, src); err != nil {
		return err
	}

	return err
}
