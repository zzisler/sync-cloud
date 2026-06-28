package updater

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func GetNewYtDlpVersion() (string, error) {
	if _, err := os.Stat(ytdlpPath()); os.IsNotExist(err) {
		return "", nil
	}
	cmd := exec.Command(ytdlpPath(), "--version")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("[yt-dlp] failed to get version: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func DownloadYtDlp(release *githubRelease) error {

	assetName := "yt-dlp"
	if runtime.GOOS == "windows" {
		assetName = "yt-dlp.exe"
	}

	var found *asset

	for _, a := range release.Assets {
		if a.Name == assetName {
			found = &a
			break
		}
	}

	if found == nil {
		return fmt.Errorf("asset %s not found in release", assetName)
	}

	data, err := DowndoadFile(found.DownloadURL, "yt-dlp")
	if err != nil {
		return err
	}

	dirPath := "./bin/yt-dlp/"

	if runtime.GOOS == "windows" {
		dirPath = filepath.Join(dirPath, "win")
	} else {
		dirPath = filepath.Join(dirPath, "linux")
	}

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s for download yt-dlp: %w", dirPath, err)
	}

	filePath := filepath.Join(dirPath, assetName)
	tempFile := filePath + ".new"

	if err := os.WriteFile(tempFile, data, 0755); err != nil {
		return fmt.Errorf("failed to write temp file yt-dlp: %w", err)
	}

	if err := os.Rename(tempFile, filePath); err != nil {
		return fmt.Errorf("failed to rename yt-dlp file: %w", err)
	}

	return nil

}

func ytdlpPath() string {
	os := runtime.GOOS
	if os == "windows" {
		return "./bin/yt-dlp/win/yt-dlp.exe"
	}
	return "./bin/yt-dlp/linux/yt-dlp"
}
