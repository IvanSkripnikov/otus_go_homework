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
	readFile, errOpen := os.Open(fromPath)
	if errOpen != nil {
		fmt.Println(errOpen)
		return nil
	}

	// проверяем, не превышает ли offset размер файла
	body, errRead := getFileBody(readFile, offset, limit)
	fmt.Println(body, errRead)

	// закрываем файл с которого читали
	defer readFile.Close()

	writeFile, errCreate := os.Create(toPath)
	if errCreate != nil {
		return errCreate
	}
	written, err := writeFile.Write([]byte(body))
	if err != nil {
		log.Panicf("failed to write: %v", err)
	}
	fmt.Println(written)
	// мы записали 1M данных !
	defer writeFile.Close() // чтобы очистить буферы ОС

	return nil
}

func getFileBody(file *os.File, offset, limit int64) (string, error) {
	var readSize, readen, counter int64

	readSize = 1024
	stopReadFlag := false

	var buf []byte
	for {
		if stopReadFlag {
			break
		}
		buf = make([]byte, readSize)
		counter++

		for readen <= readSize*counter {
			if readen < offset {
				readen++
				buf = make([]byte, readSize)
				continue
			}
			read, err := file.Read(buf[:])
			readen += int64(read)
			if err == io.EOF {
				stopReadFlag = true
				break
			}
			if err != nil {
				log.Panicf("failed to read: %v", err)
			}
		}
	}

	if offset > readen {
		return "", ErrOffsetExceedsFileSize
	}

	if limit > readen {
		limit = readen
	}

	if limit > 0 {
		buf = buf[offset : offset+limit]
	} else {
		buf = buf[offset:limit]
	}

	return string(buf), nil
}
