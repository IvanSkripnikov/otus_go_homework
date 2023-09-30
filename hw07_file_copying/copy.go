package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
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
	var nlCounterBack, nlCounterFront int64

	fileStat, _ := file.Stat()
	fileSize := fileStat.Size()
	bar := pb.StartNew(int(fileSize))

	for scanner.Scan() {
		bar.Add(len(scanner.Text() + "\n"))
		if limit > int64(len(output))-offset {
			nlCounterBack++
		}
		output += scanner.Text() + "\n"
	}

	// если была только одна строка
	if nlCounterBack == 1 {
		bar.Add(-1)
	}

	bar.Finish()
	fmt.Println("\\n")

	if offset == 0 {
		nlCounterBack--
	}

	// проверяем, не превышает ли offset размер файла
	if offset > fileSize {
		return "", ErrOffsetExceedsFileSize
	}
	output = cutOutput(output, fileSize, nlCounterFront, offset, limit)

	return output, nil
}

func cutOutput(output string, fileSize, nlCounterFront, offset, limit int64) string {
	// если limit больше размера файла - обнуляем его
	if limit > fileSize {
		limit = 0
	}

	// проверяем, не выходим ли за границы строки
	finalLength := offset + limit
	if finalLength > fileSize {
		finalLength = fileSize
	}

	if limit > 0 {
		output = output[offset+nlCounterFront : finalLength]
	} else {
		output = output[offset+nlCounterFront:]
	}

	return output
}
