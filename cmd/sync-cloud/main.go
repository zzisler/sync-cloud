package main

import (
	"log"
	"sync-cloud/internal/config"
	db "sync-cloud/internal/database"
	"sync-cloud/internal/storage"
	"sync-cloud/internal/sync"
	"sync-cloud/internal/updater"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.Connect(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	storageInstance := storage.NewStorage(database)
	if updater.CheckForUpdates(storageInstance, "yt-dlp") {
		logEntry, err := updater.Updater(storageInstance, "yt-dlp", "yt-dlp/yt-dlp")
		if err != nil {
			log.Printf("[yt-dlp] update check failed: %s", err)
		} else {
			if err := storageInstance.CreateUpdateLog(logEntry); err != nil {
				log.Printf("[yt-dlp] failed to save update log: %s", err)
			}
		}
	}

	if updater.CheckForUpdates(storageInstance, "ffmpeg") {
		logEntry, err := updater.Updater(storageInstance, "ffmpeg", "BtbN/FFmpeg-Builds")
		if err != nil {
			log.Printf("[ffmpeg] update check failed: %s", err)
		} else {
			if err := storageInstance.CreateUpdateLog(logEntry); err != nil {
				log.Printf("[ffmpeg] failed to save update log: %s", err)
			}
		}
	}

	err = sync.Run(database, cfg.ProfileURL, cfg.MusicDir, cfg.ProxyURL, cfg.DownloadTimeout)
	if err != nil {
		log.Fatal(err)
	}

}
