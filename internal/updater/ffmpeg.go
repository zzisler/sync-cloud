package updater

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync-cloud/internal/storage"
)

func GetNewFfmpegVersion(database *storage.Storage) (string, error) {
	lastLog, err := database.GetLastUpdateLog("ffmpeg")
	if err != nil {
		return "", nil
	}
	return lastLog.NewVersion, nil
}

func DownloadFfmpeg(release *githubRelease) error {

	assetName := "ffmpeg"
	if runtime.GOOS == "windows" {
		assetName = "ffmpeg.exe"
	}

	found := findFfmpegAsset(release.Assets)
	if found == nil {
		return fmt.Errorf("asset %s in not found in release", assetName)
	}

	data, err := DowndoadFile(found.DownloadURL, "ffmpeg")
	if err != nil {
		return err
	}

	dirPath := "./bin/ffmpeg/"
	if runtime.GOOS == "windows" {
		dirPath = filepath.Join(dirPath, "win")
	} else {
		dirPath = filepath.Join(dirPath, "linux")
	}

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s for download ffmpeg: %w", dirPath, err)
	}

	var extractedData []byte
	if runtime.GOOS == "windows" {
		extractedData, err = extractFfmpegFromZip(data)
	} else {
		extractedData, err = extractFfmpegFromTar(data, dirPath)
	}
	if err != nil {
		return err
	}

	filePath := filepath.Join(dirPath, assetName)
	tempFile := filePath + ".new"

	if err := os.WriteFile(tempFile, extractedData, 0755); err != nil {
		return fmt.Errorf("failed to write temp file ffmpeg: %w", err)
	}

	if err := os.Rename(tempFile, filePath); err != nil {
		return fmt.Errorf("failed to rename ffmpeg file: %w", err)
	}

	return nil

}

func extractFfmpegFromZip(data []byte) ([]byte, error) {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "ffmpeg.exe") {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("filed to open ffmpeg in zip: %w", err)
			}
			defer rc.Close()

			fileData, err := io.ReadAll(rc)
			if err != nil {
				return nil, fmt.Errorf("failed to read ffmpeg from zip: %w", err)
			}

			return fileData, nil
		}
	}

	return nil, fmt.Errorf("ffmpeg.exe not found inside zip archive")
}

func extractFfmpegFromTar(data []byte, dirPath string) ([]byte, error) {
	file, err := os.CreateTemp(dirPath, "ffmpeg-*.tar.xz")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file ffmpeg: %w", err)
	}
	defer os.Remove(file.Name())

	if _, err := file.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write archive: %w", err)
	}
	file.Close()

	cmd := exec.Command("tar", "-xJf", file.Name(), "-C", dirPath)

	if _, err = cmd.Output(); err != nil {
		return nil, fmt.Errorf("failed to extract tar.xz: %w", err)
	}

	var foundPath string

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "ffmpeg" {
			foundPath = path
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	if foundPath == "" {
		return nil, fmt.Errorf("ffmpeg binary not found after extraction")
	}

	foundData, err := os.ReadFile(foundPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read foundData file: %w", err)
	}

	return foundData, nil
}

func findFfmpegAsset(assets []asset) *asset {
	wantsOS := "linux64"
	if runtime.GOOS == "windows" {
		wantsOS = "win64"
	}
	for _, a := range assets {
		name := strings.ToLower(a.Name)
		if strings.Contains(name, wantsOS) &&
			strings.Contains(name, "gpl") &&
			!strings.Contains(name, "lgpl") &&
			!strings.Contains(name, "shared") {
			return &a
		}
	}
	return nil
}

func ffmpegPath() string {
	os := runtime.GOOS
	if os == "windows" {
		return "./bin/ffmpeg/win/ffmpeg.exe"
	}
	return "./bin/ffmpeg/linux/ffmpeg"
}
