package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
)

func getResponse(path string) ([]Item, error) {
	fileNames, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil, err
	}

	for i := range fileNames {
		fileNames[i] = strings.TrimPrefix(filepath.ToSlash(fileNames[i]), filepath.ToSlash(config.WorkingDirectory))
	}

	items := make([]Item, 0)

	for _, fileName := range fileNames {
		itemType, err := getType(filepath.Join(config.WorkingDirectory, fileName))
		if err != nil {
			log.Println(err)
			continue
		}

		if itemType == Unknown {
			continue
		}

		items = append(items, Item{itemType, Path(fileName)})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ItemType < items[j].ItemType
	})

	return items, nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	items, err := getResponse(filepath.Join(config.WorkingDirectory, r.URL.Path[1:]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = t.Execute(w, items)
}

func previewHandler(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/preview/")

	b, err := previewCache.Read(filepath.Join(config.WorkingDirectory, fileName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/avif")
	w.Header().Set("Cache-Control", "max-age=31536000")
	_, _ = w.Write(b)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/view/")

	http.ServeFile(w, r, filepath.Join(config.WorkingDirectory, fileName))
}
