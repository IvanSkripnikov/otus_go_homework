package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	inputPath  = "testdata/input.txt"
	outputPath = "/tmp/out.txt"
)

type InputData struct {
	from      string // Путь к исходному файлу
	to        string // Путь к выходному файлу
	offset    int64  // Смещение в исходном файле
	limit     int64  // Необходимое количество записываемых байтов
	reference string // Путь к эталонному файлу для сравнения
}

// Положительные сценарии копирования файлов.
func TestCopy(t *testing.T) {
	tests := []struct {
		name  string
		input InputData
	}{
		{name: "full coping", input: InputData{
			from: inputPath, to: outputPath, offset: 0, limit: 0,
			reference: "testdata/out_offset0_limit0.txt",
		}},
		{name: "coping from 0 limit 1000", input: InputData{
			from: inputPath, to: outputPath, offset: 0, limit: 1000,
			reference: "testdata/out_offset0_limit1000.txt",
		}},
		{name: "coping from 1000 limit 2000", input: InputData{
			from: inputPath, to: outputPath, offset: 1000, limit: 2000,
			reference: "testdata/out_offset1000_limit2000.txt",
		}},
		{name: "coping from 2500 limit 3050", input: InputData{
			from: inputPath, to: outputPath, offset: 2500, limit: 3050,
			reference: "testdata/out_offset2500_limit3050.txt",
		}},
		{name: "coping from -200 limit 1000", input: InputData{
			from: inputPath, to: outputPath, offset: -200, limit: 1000,
			reference: "testdata/out_offset-200_limit1000.txt",
		}},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := Copy(tc.input.from, tc.input.to, tc.input.offset, tc.input.limit)
			require.NoError(t, err)

			// Получаем контрольную сумму выходного файла
			outputSum, errOutputSum := getMd5Sum(tc.input.to)
			require.NoError(t, errOutputSum)

			// Получаем контрольную сумму эталонного файла
			referenceSum, errReferenceSum := getMd5Sum(tc.input.reference)
			require.NoError(t, errReferenceSum)

			// Проверяем соответствие между выходным и эталонным файлом по контрольной сумме
			require.Equal(t, outputSum, referenceSum, "File checksums do not match")

			// Получаем размер выходного файла
			output, outputSize, errOutputSize := getReadFileAndHimSize(tc.input.to)
			defer closeFile(output)

			require.NoError(t, errOutputSize)
			require.NotEqual(t, outputSize, 0)

			// Получаем размер эталонного файла
			reference, referenceSize, errReferenceSize := getReadFileAndHimSize(tc.input.reference)
			defer closeFile(reference)

			require.NoError(t, errReferenceSize)
			require.NotEqual(t, referenceSize, 0)

			// Проверяем соответствие между выходным и эталонным файлом по размеру файлов
			require.Equal(t, outputSize, referenceSize, "File sizes do not match")
		})
	}
}

// Тест для проверки конкретных ошибок.
func TestNegativeCopy(t *testing.T) {
	tests := []struct {
		name  string
		err   error
		input InputData
	}{
		{
			name: "from path file does not exists", err: ErrFromPathDoesNotExists,
			input: InputData{},
		},
		{
			name: "to path file does not exists", err: ErrToPathDoesNotExists,
			input: InputData{from: inputPath},
		},
		{
			name: "unsupported file", err: ErrUnsupportedFile,
			input: InputData{from: "/dev/urandom", to: outputPath},
		},
		{
			name: "offset exceeds file size", err: ErrOffsetExceedsFileSize,
			input: InputData{from: inputPath, to: outputPath, offset: 7000, limit: 0},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := Copy(tc.input.from, tc.input.to, tc.input.offset, tc.input.limit)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

// Тест на отсутствие файлов.
func TestDoesNotExistsFile(t *testing.T) {
	t.Run("does not input file", func(t *testing.T) {
		notInputPath := "testdata/not_input.txt"
		input := InputData{from: notInputPath, to: outputPath, offset: 100, limit: 2000}
		err := Copy(input.from, input.to, input.offset, input.limit)
		require.Error(t, err)
	})

	t.Run("does not output file", func(t *testing.T) {
		notOutputPath := "./test/out.txt"
		input := InputData{from: inputPath, to: notOutputPath, offset: 1000, limit: 2000}
		err := Copy(input.from, input.to, input.offset, input.limit)
		require.Error(t, err)
	})
}

// Получить контрольную сумму файла.
func getMd5Sum(filePath string) (string, error) {
	file, errFile := os.Open(filePath)
	if errFile != nil {
		return "", errFile
	}

	fileSum := md5.New()
	_, errSum := io.Copy(fileSum, file)
	if errSum != nil {
		return "", errSum
	}

	return fmt.Sprintf("%X", fileSum.Sum(nil)), nil
}
