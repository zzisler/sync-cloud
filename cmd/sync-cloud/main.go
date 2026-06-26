package main

import (
	"log"
	"sync-cloud/internal/config"
	db "sync-cloud/internal/database"
	"sync-cloud/internal/sync"
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

	err = sync.Run(database, cfg.ProfileURL, cfg.MusicDir, cfg.ProxyURL)
	if err != nil {
		log.Fatal(err)
	}

}
