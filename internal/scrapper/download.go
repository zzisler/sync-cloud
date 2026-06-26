package scrapper

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var filenameSanitizer = strings.NewReplacer(
	"/", "_",
	"\\", "_",
	":", "_",
	"*", "_",
	"?", "_",
	"\"", "_",
	"<", "_",
	">", "_",
	"|", "_",
)

func DownloadAndTag(track *Track, musicDir, proxyURL string) (string, error) {

	client, err := newHttpClient(proxyURL)
	if err != nil {
		return "", fmt.Errorf("failed to create http client: %w", err)
	}

	outputPath := filepath.Join(musicDir, buildFilename(track))

	dataCover, err := DownloadCover(track.CoverURL, client)
	if err != nil {
		return "", err
	}

	tempFile, err := os.CreateTemp("", "cover-*.jpg")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(dataCover); err != nil {
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}

	tempFile.Close()

	ytDlpCmd := exec.Command(ytdlpPath(), "-o", "-", track.URL)

	ffmpegCmd := exec.Command(ffmpegPath(),
		"-i", "-",
		"-i", tempFile.Name(),
		"-map", "0:a", "-map", "1:0",
		"-c:a", "libmp3lame", "-b:a", "192k",
		"-c:v", "copy",
		"-id3v2_version", "3",
		"-metadata", "title="+track.Title,
		"-metadata", "artist="+track.Artist,
		"-disposition:1", "attached_pic",
		outputPath,
	)

	pipe, err := ytDlpCmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	ffmpegCmd.Stdin = pipe

	if err := ffmpegCmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start ffmpeg: %w", err)
	}
	if err := ytDlpCmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start yt-dlp: %w", err)
	}

	if err := ytDlpCmd.Wait(); err != nil {
		return "", fmt.Errorf("yt-dlp failed: %w", err)
	}
	if err := ffmpegCmd.Wait(); err != nil {
		return "", fmt.Errorf("ffmpeg failed: %w", err)
	}

	return outputPath, nil
}

func newHttpClient(proxyURL string) (*http.Client, error) {
	if proxyURL == "" {
		return http.DefaultClient, nil
	}

	parsedProxy, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy URL: %w", err)
	}

	transtort := &http.Transport{
		Proxy: http.ProxyURL(parsedProxy),
	}

	return &http.Client{Transport: transtort}, nil

}

func DownloadCover(url string, client *http.Client) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download cover: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	return data, nil
}

func buildFilename(track *Track) string {
	return fmt.Sprintf("%s - %s.mp3", filenameSanitizer.Replace(track.Title), filenameSanitizer.Replace(track.Artist))
}
