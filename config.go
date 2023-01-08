package main

import (
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	WorkingDirectory string `env-default:"."`
	Crf              uint64 `env-default:"40"`
	ProcessCount     uint64 `env-default:"4"`
}

func loadConfig(path string) (*Config, error) {
	var config Config

	err := cleanenv.ReadConfig("gallery.toml", &config)
	if err != nil {
		return nil, err
	}

	workingDirectory, err := filepath.Abs(config.WorkingDirectory)
	if err != nil {
		return nil, err
	}
	config.WorkingDirectory = workingDirectory

	return &config, nil
}
