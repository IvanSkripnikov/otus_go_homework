package main

import (
	"errors"
	"github.com/cheggaaa/pb/v3"
	"io"
	"os"
)

var (
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOpenFile              = errors.New("can't open file")
	ErrCreateFile            = errors.New("can't create file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
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

	if limit > fileSize {
		limit = fileSize
	}

	bufferSize := 1024

	// Настраиваем прогресс бар
	progressCounts := getProgressBarLimit(fileSize, offset)

	bar := pb.StartNew(progressCounts)
	bar.Set(pb.Bytes, true)
	defer bar.Finish()

	var (
		completeHandleCount int64
		readLen             int64
		skippedCount        int64
	)

	if offset > 0 {
		writeFile, errWrite := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if errWrite != nil {
			return ErrCreateFile
		}

		for skippedCount <= offset {
			if offset < int64(bufferSize) {
				readLen = offset
			} else {
				readLen = int64(bufferSize)
			}
			if _, err := io.CopyN(writeFile, readFile, readLen); err != nil {
				break
			}

			skippedCount += readLen
		}

		writeFile.Close()
	}

	writeFile, errWrite := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if errWrite != nil {
		return ErrCreateFile
	}

	for completeHandleCount < limit {
		if limit < int64(bufferSize) {
			readLen = limit
		} else {
			readLen = int64(bufferSize)
		}
		if _, err := io.CopyN(writeFile, readFile, readLen); err != nil {
			break
		}

		completeHandleCount += readLen
		//fmt.Println(completeHandleCount, limit, fileSize, readLen)
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
