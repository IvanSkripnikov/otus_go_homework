package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOpenFile              = errors.New("can't open file")
	ErrCreateFile            = errors.New("can't create file")
)

func Copy(fromPath, toPath string, offset, limit int64) error { // мы заранее знаем сколько хотим прочитать
	// открываем файл
	readFile, errOpen := os.Open(fromPath)
	if errOpen != nil {
		return ErrOpenFile
	}

	body, errRead := getFileBody(readFile, offset, limit)
	if errRead != nil {
		return errRead
	}
	fmt.Println(body)
	writeFile, errCreate := os.Create(toPath)
	if errCreate != nil {
		return ErrCreateFile
	}
	written, err := writeFile.Write([]byte(body))
	if err != nil {
		log.Panicf("failed to write: %v", err)
	}
	fmt.Println(written)

	// закрываем файл с которыми работали
	defer func() {
		writeFile.Close()
		readFile.Close()
	}()

	return nil
}

func getFileBody(file *os.File, offset, limit int64) (string, error) {
	var readSize, readen int64
	readSize = 1024
	var buf []byte

	for {
		fmt.Println("cycle", readen, readSize, offset, limit)
		read, err := file.Read(buf[readen:])
		readen += int64(read)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panicf("failed to read: %v", err)
		}
		if readen < offset {
			buf = make([]byte, readSize)
		}
	}

	// проверяем, не превышает ли offset размер файла
	if offset > readen {
		return "", ErrOffsetExceedsFileSize
	}

	if limit > readen {
		limit = readen
	}

	if limit > 0 {
		buf = buf[offset : offset+limit]
	} else {
		buf = buf[offset:]
	}

	return string(buf), nil
}
