package main

import (
	"bufio"
	"errors"
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

	writeFile, errCreate := os.Create(toPath)
	if errCreate != nil {
		return ErrCreateFile
	}
	_, err := writeFile.Write([]byte(body))
	if err != nil {
		log.Panicf("failed to write: %v", err)
	}

	// закрываем файл с которыми работали
	defer func() {
		writeFile.Close()
		readFile.Close()
	}()

	return nil
}

func getFileBody(file *os.File, offset, limit int64) (string, error) {
	output := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		output += scanner.Text()
	}

	// проверяем, не превышает ли offset размер файла
	if offset > int64(len(output)) {
		return "", ErrOffsetExceedsFileSize
	}

	if limit > 0 {
		output = output[offset : offset+limit]
	} else {
		output = output[offset:]
	}

	return output, nil
}
