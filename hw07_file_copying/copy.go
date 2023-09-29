package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOpenFile              = errors.New("can't open file")
	ErrCreateFile            = errors.New("can't create file")
)

func Copy(fromPath, toPath string, offset, limit int64) error { // мы заранее знаем сколько хотим прочитать
	extension := filepath.Ext(fromPath)
	if extension != ".txt" {
		return ErrUnsupportedFile
	}
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
	var nlCounterBack, nlCounterFront, lenOutput int64

	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		if limit > int64(len(output))-offset {
			//		fmt.Println(nlCounterBack, limit, int64(len(output))-offset)
			nlCounterBack++

		}
		if int64(len(output)) < offset {
			nlCounterFront++
		}
		output += scanner.Text() + "\r\n"
	}

	if offset == 0 {
		nlCounterBack--
	}

	lenOutput = int64(len(output))

	// проверяем, не превышает ли offset размер файла
	if offset > lenOutput {
		return "", ErrOffsetExceedsFileSize
	}

	// если limit больше размера файла - обнуляем его
	if limit > lenOutput {
		limit = 0
	}

	finalLength := offset + limit + nlCounterBack
	if finalLength > lenOutput {
		finalLength = lenOutput
		nlCounterFront--
	}

	if limit > 0 {
		output = output[offset+nlCounterFront : finalLength]
	} else {
		output = output[offset+nlCounterFront:]
	}

	return output, nil
}
