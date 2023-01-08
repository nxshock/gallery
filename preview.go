package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/nxshock/zkv"
)

type PreviewCache struct {
	store *zkv.Store
}

func NewPreviewCache(filePath string) *PreviewCache {
	store, err := zkv.Open(filePath)
	if err != nil {
		log.Fatalln(err)
	}

	reviewCache := &PreviewCache{store}

	return reviewCache
}

func (pc *PreviewCache) Add(filePath string) ([]byte, error) {
	defer func() {
		<-semaphore
	}()
	semaphore <- struct{}{}

	tempFileName, err := TempFileName("preview*.avif")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFileName)

	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	var cmd *exec.Cmd
	if stat.IsDir() { // https://trac.ffmpeg.org/wiki/Create%20a%20mosaic%20out%20of%20several%20input%20videos
		fileNames, err := getRandomFiles(filePath, 4)
		if err != nil {
			return nil, err
		}

		cmd = exec.Command("ffmpeg",
			"-i", fileNames[0],
			"-i", fileNames[1],
			"-i", fileNames[2],
			"-i", fileNames[3],
			"-filter_complex", "nullsrc=size=240x240 [base];[0:v] setpts=PTS-STARTPTS, scale=120x120:force_original_aspect_ratio=increase [upperleft];[1:v] setpts=PTS-STARTPTS, scale=120x120:force_original_aspect_ratio=increase [upperright];[2:v] setpts=PTS-STARTPTS, scale=120x120:force_original_aspect_ratio=increase [lowerleft];[3:v] setpts=PTS-STARTPTS, scale=120x120:force_original_aspect_ratio=increase [lowerright];[base][upperleft] overlay=shortest=1 [tmp1];[tmp1][upperright] overlay=shortest=1:x=120 [tmp2];[tmp2][lowerleft] overlay=shortest=1:y=120 [tmp3];[tmp3][lowerright] overlay=shortest=1:x=120:y=120",
			"-frames:v", "1",
			"-crf", strconv.FormatUint(config.Crf, 10),
			"-f", "avif",
			tempFileName)
	} else {
		cmd = exec.Command("ffmpeg.exe",
			"-i", filepath.FromSlash(filePath),
			"-vf", "scale=240:240:force_original_aspect_ratio=increase,crop=240:240:exact=1",
			"-frames:v", "1",
			"-crf", "40",
			"-f", "avif",
			tempFileName)
	}

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(tempFileName)
	if err != nil {
		return nil, err
	}

	err = pc.store.Set(filePath, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (pc *PreviewCache) Read(filePath string) ([]byte, error) {
	var b []byte

	err := pc.store.Get(filePath, &b)
	if err != nil {
		return pc.Add(filePath)
	}

	return b, nil
}

func (pc *PreviewCache) Save() error {
	return pc.store.Close()
}
