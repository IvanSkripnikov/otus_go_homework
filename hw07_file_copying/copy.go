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
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFromPathDoesNotExists = errors.New("from path a file does not exist")
	ErrToPathDoesNotExists   = errors.New("to path a file does not exist")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	var (
		errOffset   error
		hasEndWrite bool
		writeOffset int64
	)

	errPaths := checkExistsPathsOfFiles(fromPath, toPath)
	if errPaths != nil {
		return errPaths
	}

	input, inputFileSize, errInput := getReadFileAndHimSize(fromPath)
	if errInput != nil {
		return errInput
	}

	// Закрываем входной файл
	defer closeFile(input)

	output, errOutput := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if errOutput != nil {
		errMessage := fmt.Sprintf("failed write to output file, error: %v", errOutput)
		return errors.New(errMessage)
	}

	// Закрываем выходной файл
	defer closeFile(output)

	bufferSize := 1024
	buffer := make([]byte, bufferSize)
	clearBuffer := make([]byte, bufferSize)

	// Проверка значения для смещения в исходном файле, и при необходимости его корректировка
	offset, errOffset = correctOffsetForNegativeValue(offset, inputFileSize)
	if errOffset != nil {
		return errOffset
	}

	// Настраиваем прогресс бар
	progressCounts := getProgressCounts(inputFileSize, offset)

	bar := pb.StartNew(progressCounts)
	bar.Set(pb.Bytes, true)
	defer bar.Finish()

	for offset < inputFileSize {
		read, errRead := input.ReadAt(buffer, offset)
		if errRead != nil && errRead != io.EOF {
			errMessage := fmt.Sprintf("error reading from input file, error: %v", errRead)
			return errors.New(errMessage)
		}

		switch {
		// Если прочитали меньше чем ожидалось, то уменьшаем размер буфера
		case read < bufferSize:
			buffer = buffer[:read]
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

		written, errWrite := output.WriteAt(buffer, writeOffset)
		if errWrite != nil {
			errRemove := os.Remove(toPath)
			if errRemove != nil {
				log.Println(errRemove)
			}

			errMessage := fmt.Sprintf("error writing to output file, error: %v", errWrite)
			return errors.New(errMessage)
		}

		// Добавляем количество записанных байт в прогрессбар
		bar.Add(written)

		// Очищаем буфер с данными
		copy(buffer, clearBuffer)

		// Смещаем позицию в файле
		offset += int64(read)
		writeOffset += int64(read)

		// Если файл дочитан до конца, или установлен флаг
		// что запрошенный объем данных уже был записан, то выходим из цикла
		if errRead == io.EOF || hasEndWrite {
			break
		}
	}

	return nil
}

// Проверка наличия путей к файлам.
func checkExistsPathsOfFiles(fromPath, toPath string) error {
	if len(fromPath) == 0 {
		return ErrFromPathDoesNotExists
	}

	if len(toPath) == 0 {
		return ErrToPathDoesNotExists
	}

	return nil
}

// Открыть файл для чтения и получить информацию о его размере.
func getReadFileAndHimSize(filePath string) (*os.File, int64, error) {
	file, errFile := os.Open(filePath)
	if errFile != nil {
		errMessage := fmt.Sprintf("failed to read file, error: %v", errFile)
		return nil, 0, errors.New(errMessage)
	}

	info, errInfo := file.Stat()
	if errInfo != nil {
		errMessage := fmt.Sprintf("failed read info meta data from file, error: %v", errInfo)
		return nil, 0, errors.New(errMessage)
	}

	inputFileSize := info.Size()
	if inputFileSize <= 0 {
		return nil, 0, ErrUnsupportedFile
	}

	return file, inputFileSize, nil
}

// Закрыть файл.
func closeFile(file *os.File) {
	if file != nil {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}
}

// Корректировка смещения если было введено отрицательное значение.
func correctOffsetForNegativeValue(offset, inputFileSize int64) (int64, error) {
	var negativeOffset int64
	if offset < 0 {
		negativeOffset = offset * (-1)
		offset = inputFileSize - negativeOffset
	}

	if offset > inputFileSize || negativeOffset > inputFileSize {
		return 0, ErrOffsetExceedsFileSize
	}

	return offset, nil
}

// Получить количество записывемых байт для прогресс бара.
func getProgressCounts(inputFileSize, offset int64) int {
	progressCounts := int(inputFileSize - offset)

	if limit > 0 && limit < (inputFileSize-offset) {
		progressCounts = int(limit)
	}

	return progressCounts
}
