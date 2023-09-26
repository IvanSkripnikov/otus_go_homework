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
)

func Copy(fromPath, toPath string, offset, limit int64) error { // мы заранее знаем сколько хотим прочитать
	// открываем файл
	file, errOpen := os.Open(fromPath)
	if errOpen != nil {
		fmt.Println(errOpen)
		return nil
	}

	// проверяем, не превышает ли offset размер файла
	if isOffsetExceedsFileSize(file, offset) {
		return ErrOffsetExceedsFileSize
	}

	b, err := os.ReadFile(fromPath)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(b))

	file, errCreate := os.Create(toPath)

	if errCreate != nil {
		return errCreate
	}
	written, err := file.Write(b)
	if err != nil {
		log.Panicf("failed to write: %v", err)
	}
	fmt.Println(written)
	// мы записали 1M данных !
	defer file.Close() // чтобы очистить буферы ОС

	return nil
}

func isOffsetExceedsFileSize(file *os.File, offset int64) bool {
	var readSize int64
	var readen int64
	readSize = 1024

	buf := make([]byte, readSize)
	for readen < readSize {
		read, err := file.Read(buf[readen:])

		readen += int64(read)
		if err == io.EOF {
			// что если не дочитали ?
			break
		}
		if err != nil {
			log.Panicf("failed to read: %v", err)
		}
	}
	fmt.Println("output", string(buf))
	return offset > readen
}
