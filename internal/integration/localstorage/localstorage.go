package localstorage

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

type LocalStorageContract interface {
}

type localstorage struct {
}

func NewLocalStorageIntegration() LocalStorageContract {
	return &localstorage{}
}

func (l *localstorage) Save(base64String, filename, path string) (fullpath string, err error) {
	// Trim base64 prefix if present (e.g., "data:image/png;base64,")
	if idx := strings.Index(base64String, ","); idx != -1 {
		base64String = base64String[idx+1:]
	}

	// Decode base64 string to byte slice
	data, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		log.Error().Err(err).Msg("failed to decode base64 string")
		return "", fmt.Errorf("failed to decode base64 string: %w", err)
	}

	// Save file to local storage
	fullpath = fmt.Sprintf("%s/%s", path, filename)
	if err := l.saveFile(fullpath, data); err != nil {
		log.Error().Err(err).Msg("failed to save file")
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return fullpath, nil
}

func (l *localstorage) saveFile(fullpath string, data []byte) error {
	path := strings.Split(fullpath, "/")
	dir := strings.Join(path[:len(path)-1], "/")

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Error().Err(err).Msg("failed to create directory")
		return err
	}

	file, err := os.Create(fullpath)
	if err != nil {
		log.Error().Err(err).Msg("failed to create file")
		return err
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		log.Error().Err(err).Msg("failed to write data to file")
		return err
	}

	return nil
}
