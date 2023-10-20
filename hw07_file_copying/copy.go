package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrOffsetExceedsFileSize       = errors.New("offset exceeds file size")
	ErrOpenFile                    = errors.New("can't open file")
	ErrCreateFile                  = errors.New("can't create file")
	bufferSize               int64 = 1024
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

	// Настраиваем прогресс бар
	progressCounts := getProgressBarLimit(fileSize, offset)

	bar := pb.StartNew(progressCounts)
	bar.Set(pb.Bytes, true)
	defer bar.Finish()

	var (
		completeHandleCount int64
		readLen             int64
	)

	if offset > 0 {
		rewind := rewindOffset(*readFile, toPath, bar)
		if rewind != nil {
			return rewind
		}
	}

	writeFile, errWrite := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if errWrite != nil {
		return ErrCreateFile
	}

	if limit < bufferSize {
		readLen = limit
	} else {
		readLen = bufferSize
	}
	for completeHandleCount < limit {
		if _, err := io.CopyN(writeFile, readFile, readLen); err != nil {
			break
		}

		completeHandleCount += readLen

		// инкерментим прогрессбар
		bar.Add(int(readLen))
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

func rewindOffset(readFile os.File, toPath string, bar *pb.ProgressBar) error {
	var (
		readLen      int64
		skippedCount int64
	)

	writeFile, errWrite := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if errWrite != nil {
		return ErrCreateFile
	}

	if offset < bufferSize {
		readLen = offset
	} else {
		readLen = bufferSize
	}

	for skippedCount < offset {
		if skippedCount+readLen > offset {
			readLen = offset - skippedCount
		}

		if _, err := io.CopyN(writeFile, &readFile, readLen); err != nil {
			break
		}

		skippedCount += readLen

		// инкерментим прогрессбар
		bar.Add(int(readLen))
	}
	writeFile.Close()

	return nil
}
