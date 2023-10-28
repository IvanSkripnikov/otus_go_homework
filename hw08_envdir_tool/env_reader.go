package main

import (
	"bytes"
	"log"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	nl := []byte("\n")
	zero := []byte("\x00")

	// Читаем все файлы из каталога
	files, errDir := os.ReadDir(dir)
	if errDir != nil {
		return nil, errDir
	}

	// cоздаём массив переменных среды
	envs := make(Environment, len(files))

	for _, file := range files {
		if canSkipObject(file, &envs) {
			continue
		}
		envVarName := file.Name()

		// получаем данные
		data, errData := os.ReadFile(dir + "/" + envVarName)
		if errData != nil {
			log.Printf("Can't read file: %s, error: %v \n", envVarName, errData)
			continue
		}

		// разбиваем данные на массив строк
		lines := bytes.Split(data, nl)

		// в первой строке заменяем нули на перенос строки
		firstLine := bytes.ReplaceAll(lines[0], zero, nl)

		// убираем пробельные символы с краёв строки
		envValue := trim(firstLine)

		// если такая переменная установлена - удаляем
		checkExistsEnv(envVarName)

		envs[envVarName] = EnvValue{Value: envValue}
	}

	return envs, nil
}

func removeEnvVar(key string) {
	err := os.Unsetenv(key)
	if err != nil {
		log.Printf("Can't delete environment variable: %s, error: %v \n", key, err)
	}
}

func canSkipObject(file os.DirEntry, envs *Environment) bool {
	result := false
	envVarName := file.Name()

	// если текущий объект директория или название содержит =, пропускаем
	if file.IsDir() || strings.Contains(envVarName, "=") {
		result = true
	}

	// если не можем получить данных по файлу - пропускаем
	info, errInfo := file.Info()
	if errInfo != nil {
		log.Printf("Can't get file info: %s, error: %v \n", envVarName, errInfo)
		result = true
	}

	// если файл пуст - пропускаем, и удаляем переменную
	if info.Size() == 0 {
		(*envs)[envVarName] = EnvValue{NeedRemove: true}
		removeEnvVar(envVarName)
		result = true
	}

	return result
}

func checkExistsEnv(name string) {
	_, ok := os.LookupEnv(name)
	if ok {
		removeEnvVar(name)
	}
}

func trim(value []byte) string {
	return strings.TrimRight(string(value), " \t")
}
