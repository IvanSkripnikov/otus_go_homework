package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
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
	progressCounts := getProgressCounts(fileSize, offset)

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

		written, errWrite := writeFile.WriteAt(buffer, writeOffset)
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

	// закрываем файл с которыми работали
	defer func() {
		writeFile.Close()
		readFile.Close()
	}()

	return nil
}

/*
	func getFileBody(file *os.File, offset, limit int64) (string, error) {
		output := ""
		scanner := bufio.NewScanner(file)
		var nlCounterBack int64

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

		// проверяем, не превышает ли offset размер файла
		if offset > fileSize {
			return "", ErrOffsetExceedsFileSize
		}
		output = cutOutput(output, fileSize, offset, limit)

		return output, nil
	}

	func cutOutput(output string, fileSize, offset, limit int64) string {
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
			output = output[offset:finalLength]
		} else {
			output = output[offset:]
		}

		return output
	}
*/
func getProgressCounts(inputFileSize, offset int64) int {
	progressCounts := int(inputFileSize - offset)

	if limit > 0 && limit < (inputFileSize-offset) {
		progressCounts = int(limit)
	}

	return progressCounts
}
