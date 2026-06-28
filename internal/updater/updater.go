package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync-cloud/internal/storage"
	"time"
)

type githubRelease struct {
	TagName string  `json:"tag_name"`
	Name    string  `json:"name"`
	Assets  []asset `json:"assets"`
}

type asset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
}

func CheckForUpdates(database *storage.Storage, tool string) bool {
	lastLog, err := database.GetLastUpdateLog(tool)
	if err != nil {
		return true
	}
	return time.Since(lastLog.CheckedAt) >= 7*24*time.Hour
}

func Updater(database *storage.Storage, toolName, repo string) (*storage.UpdateLog, error) {

	var currentVersion string
	var err error

	switch toolName {
	case "yt-dlp":
		currentVersion, err = GetNewYtDlpVersion()
	case "ffmpeg":
		currentVersion, err = GetNewFfmpegVersion(database)
	default:
		return nil, fmt.Errorf("unknown tool: %s", toolName)
	}
	if err != nil {
		return nil, err
	}

	release, err := gitUpdate(repo)
	if err != nil {
		return nil, err
	}

	if toolName == "ffmpeg" {
		fields := strings.Fields(release.Name)  // [ "Latest", "Auto-Build", "(2026-06-27", "13:21)"]
		version := strings.Trim(fields[2], "(") // "2026-06-27"
		release.TagName = version
	}

	logEntry := &storage.UpdateLog{
		CheckedAt:  time.Now(),
		Tool:       toolName,
		OldVersion: currentVersion,
		NewVersion: release.TagName,
		Updated:    false,
	}

	if currentVersion == release.TagName {
		return logEntry, nil
	}

	switch toolName {
	case "yt-dlp":
		err = DownloadYtDlp(release)
	case "ffmpeg":
		err = DownloadFfmpeg(release)
	}
	if err != nil {
		return nil, err
	}

	logEntry.Updated = true
	return logEntry, nil

}

func gitUpdate(repo string) (*githubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("[%s] failed to get latest version: %w", repo, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[%s] failed to read body latest version: %w", repo, err)
	}

	var release githubRelease

	if err := json.Unmarshal(data, &release); err != nil {
		return nil, fmt.Errorf("[%s] failed to unmarshal json: %w", repo, err)
	}

	return &release, nil
}

func DowndoadFile(url, toolName string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("[%s] failed download file: %w", toolName, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[%s] failed to read body for download: %w", toolName, err)
	}

	return data, nil
}
