package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
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

	fileStat, _ := readFile.Stat()
	fileSize := fileStat.Size()

	// проверяем, не превышает ли offset размер файла
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	bufferSize := 1024
	buffer := make([]byte, bufferSize)
	clearBuffer := make([]byte, bufferSize)

	// Настраиваем прогресс бар
	progressCounts := getProgressBarLimit(fileSize, offset)

	bar := pb.StartNew(progressCounts)
	bar.Set(pb.Bytes, true)
	defer bar.Finish()

	writeFile, errWrite := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if errWrite != nil {
		return ErrCreateFile
	}

	var writeOffset int64
	hasEndWrite := false

	for offset < fileSize {
		read, errRead := readFile.ReadAt(buffer, offset)
		if errRead != nil && errRead != io.EOF {
			return errRead
		}

		switch {
		// Если прочитали меньше чем ожидалось, то уменьшаем размер буфера
		case read < bufferSize:
			buffer = buffer[:read]
			// buffer = buffer[offset : offset+limit]
			hasEndWrite = true

		// Если заданый лимит меньше чем прочитаная часть данных, то уменьшаем размер буфера
		case limit > 0 && limit < int64(read):
			buffer = buffer[:limit]
			hasEndWrite = true

		// Если текущий проход записи данных, превышает заданый лимит, то уменьшаем размер буфера
		case limit > 0 && limit < writeOffset+int64(bufferSize):
			sliceOffset := int64(math.Abs(float64(writeOffset - limit)))
			buffer = buffer[:sliceOffset]
			hasEndWrite = true
		}

		written, errWrite := writeFile.WriteAt(buffer, writeOffset)
		if errWrite != nil {
			errRemove := os.Remove(toPath)
			if errRemove != nil {
				log.Println(errRemove)
			}

			errMessage := fmt.Sprintf("error writing to output file, error: %v", errWrite)
			return errors.New(errMessage)
		}

		// инкерментим прогрессбар
		bar.Add(written)

		// очищаем буфер
		copy(buffer, clearBuffer)

		// перемещаем позицию для следующего чтения в файле
		offset += int64(read)
		writeOffset += int64(read)

		// проверяем условия выхода из цикла записи
		if errRead == io.EOF || hasEndWrite {
			break
		}
	}

	// закрываем файл с которыми работали
	defer func() {
		writeFile.Close()
		readFile.Close()
	}()

	return nil
}

func getProgressBarLimit(inputFileSize, offset int64) int {
	progressCounts := int(inputFileSize - offset)

	if limit > 0 && limit < (inputFileSize-offset) {
		progressCounts = int(limit)
	}

	return progressCounts
}
