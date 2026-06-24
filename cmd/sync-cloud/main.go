package main

import (
	"log"
	db "sync-cloud/internal/database"
	"sync-cloud/internal/sync"
)

func main() {

	database, err := db.Connect("likes.db")
	if err != nil {
		log.Fatal(err)
	}

	err = sync.Run(database, "https://soundcloud.com/zzisler-sc/likes")
	if err != nil {
		log.Fatal(err)
	}

}
