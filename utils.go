package main

import (
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var errPatternHasSeparator = errors.New("pattern contains path separator")

func TempFileName(pattern string) (string, error) {
	dir := os.TempDir()

	prefix, suffix, err := prefixAndSuffix(pattern)
	if err != nil {
		return "", &os.PathError{Op: "createtemp", Path: pattern, Err: err}
	}
	prefix = joinPath(dir, prefix)

	try := 0
	for {
		name := prefix + nextRandom() + suffix
		if exists, _ := IsFileExists(name); exists {
			if try++; try < 10000 {
				continue
			}
			return "", &os.PathError{Op: "createtemp", Path: prefix + "*" + suffix, Err: os.ErrExist}
		}
		return name, err
	}
}

func IsFileExists(filePath string) (bool, error) {
	if _, err := os.Stat(filePath); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil

	} else {
		return false, err
	}
}

func prefixAndSuffix(pattern string) (prefix, suffix string, err error) {
	for i := 0; i < len(pattern); i++ {
		if os.IsPathSeparator(pattern[i]) {
			return "", "", errPatternHasSeparator
		}
	}
	if pos := lastIndex(pattern, '*'); pos != -1 {
		prefix, suffix = pattern[:pos], pattern[pos+1:]
	} else {
		prefix = pattern
	}
	return prefix, suffix, nil
}

func joinPath(dir, name string) string {
	if len(dir) > 0 && os.IsPathSeparator(dir[len(dir)-1]) {
		return dir + name
	}
	return dir + string(os.PathSeparator) + name
}

func nextRandom() string {

	return strconv.Itoa(int(rand.Uint64()))
}

func lastIndex(s string, sep byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == sep {
			return i
		}
	}
	return -1
}

func getRandomFiles(path string, numberOfFiles int) ([]string, error) { // TODO: оптимизировать алгоритм и предусмотреть папки с менее 4 шт. файлов
	fileNames, err := filepath.Glob(path + "/*")
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	for i := len(fileNames) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		fileNames[i], fileNames[j] = fileNames[j], fileNames[i]
	}

	for len(fileNames) < 4 {
		fileNames = append(fileNames, "black.png")
	}

	return fileNames[:4], nil
}
