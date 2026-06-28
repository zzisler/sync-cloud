package scrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strconv"
)

func GetMeta(url, proxyURL string) (*Track, error) {

	// переменная с структурой метаданных
	var meta metadata

	// путь до бинарника yt-dlp
	ytDlp := ytdlpPath()

	// генерируем случайное число от 1.5 до 3.5 для запроса yt-dlp
	randomInt := fmt.Sprintf("%.2f", 1.5+rand.Float64()*(3.5-1.5))

	log.Printf("%ss | fetching metadata: %s", randomInt, url)

	args := []string{"--sleep-requests", randomInt, "--print", "%(.{id,title,artist,uploader,thumbnails})j"}
	if proxyURL != "" {
		args = append(args, "--proxy", proxyURL)
	}
	args = append(args, url)

	// запускаем бинарник и выводим нужные метаданные о треке по ссылке
	cmd := exec.Command(ytDlp, args...)

	// тут я не помню что было, но вроде мы создаем переменную для байтов и пишем вывод сюда
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// закрываем вывод и проверяем на ошибку
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp error: %s", stderr.String())
	}

	// конвертируем json с метаданными трека в переменную meta созданную в начале
	err = json.Unmarshal(out, &meta)
	if err != nil {
		return nil, err
	}

	// тут берем из массива одну конкретную ссылку на фотку оригинал
	coverUrl := findCover(&meta)
	if coverUrl == "" {
		log.Printf("no cover found for url: %s", url)
		return nil, fmt.Errorf("no cover found")
	}

	// переводим строку id трека в int64
	id, _ := strconv.ParseInt(meta.ID, 10, 64)

	// возвращаем данные о треке
	return &Track{
		ID:       id,
		Title:    meta.Title,
		Artist:   resolveArtist(&meta),
		CoverURL: coverUrl,
		URL:      url,
	}, nil

}

func findCover(meta *metadata) string {
	for i := 0; i < len(meta.Thumbnails); i++ {
		if meta.Thumbnails[i].ID == "original" {
			return meta.Thumbnails[i].URL
		}
	}
	return ""
}

func resolveArtist(meta *metadata) string {
	if meta.Artist == "" {
		return meta.Uploader
	}
	return meta.Artist
}
