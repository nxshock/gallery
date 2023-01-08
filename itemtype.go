package main

import (
	"os"
	"path/filepath"
	"strings"
)

type ItemType int

const (
	Unknown ItemType = iota
	Directory
	Picture
	Video
)

type Path string

type Item struct {
	ItemType
	Path Path
}

func (p *Path) Base() string {
	return filepath.Base(string(*p))
}

func getType(path string) (ItemType, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return Unknown, err
	}

	if stat.IsDir() {
		return Directory, nil
	}

	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg", ".png", ".avif", ".bmp", ".gif":
		return Picture, nil
	case ".mov", ".mp4", ".avi", ".mkv", ".3gp", ".webm":
		return Video, nil
	}

	return Unknown, nil
}
