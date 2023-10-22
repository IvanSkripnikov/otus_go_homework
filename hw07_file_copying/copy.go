package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
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

	if limit > fileSize || limit == 0 {
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
		// skippedCount        int64
	)

	if offset > 0 {
		/*
			writeFile, errWrite := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
			if errWrite != nil {
				return ErrCreateFile
			}

			if offset < int64(bufferSize) {
				readLen = offset
			} else {
				readLen = int64(bufferSize)
			}
			for skippedCount < offset {
				if skippedCount+readLen > offset {
					readLen = offset - skippedCount
				}
				io.CopyN(writeFile, readFile, readLen)

				skippedCount += readLen

				// инкерментим прогрессбар
				bar.Add(int(readLen))
			}
			writeFile.Close()
		*/
		_, err := readFile.Seek(offset, int(limit))
		if err != nil {
			return err
		}
	}

	writeFile, errWrite := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if errWrite != nil {
		return ErrCreateFile
	}

	if limit < int64(bufferSize) {
		readLen = limit
	} else {
		readLen = int64(bufferSize)
	}
	for completeHandleCount < limit {
		_, err := io.CopyN(writeFile, readFile, readLen)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}

		completeHandleCount += readLen

		// инкерментим прогрессбар
		bar.Add(int(readLen))
	}

	bar.Add(int(bar.Total() - bar.Current()))

	// закрываем файл с которыми работали
	defer func() {
		err := writeFile.Close()
		if err != nil {
			return
		}
		err = readFile.Close()
		if err != nil {
			return
		}
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
