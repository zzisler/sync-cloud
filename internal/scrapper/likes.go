package scrapper

import (
	"bufio"
	"encoding/json"
	"log"
	"os/exec"
	"runtime"
)

func FetchLikes(profileLink, proxyURL string) ([]LikeEntry, error) {

	ytDlp := ytdlpPath()

	args := []string{"--flat-playlist", "-j"}
	if proxyURL != "" {
		args = append(args, "--proxy", proxyURL)
	}
	args = append(args, profileLink)

	cmd := exec.Command(ytDlp, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("error cmd.StdoutPipe: %s", err)
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		log.Printf("error cmd.Start: %s", err)
		return nil, err
	}

	var tracks []LikeEntry

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {

		var track LikeEntry

		err := json.Unmarshal([]byte(scanner.Text()), &track)
		if err != nil {
			log.Printf("error unmarshal track, err: %s", err)
			continue
		}

		tracks = append(tracks, track)

	}

	if err := cmd.Wait(); err != nil {
		log.Printf("error cmd.Wait: %s", err)
		return nil, err
	}

	return tracks, nil

}

func ytdlpPath() string {
	os := runtime.GOOS
	if os == "windows" {
		return "./bin/yt-dlp/win/yt-dlp.exe"
	}
	return "./bin/yt-dlp/linux/yt-dlp"
}

func ffmpegPath() string {
	os := runtime.GOOS
	if os == "windows" {
		return "./bin/ffmpeg/win/ffmpeg.exe"
	}
	return "./bin/ffmpeg/linux/ffmpeg"
}
