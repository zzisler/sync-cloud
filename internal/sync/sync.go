package sync

import (
	"log"
	"sync-cloud/internal/scrapper"
	"sync-cloud/internal/storage"

	"gorm.io/gorm"
)

func Run(db *gorm.DB, profileLink string) error {

	log.Println("=== sync-cloud: starting sync ===")

	database := storage.NewStorage(db)

	log.Println("fetching likes from SoundCloud...")

	likesList, err := scrapper.FetchLikes(profileLink)
	if err != nil {
		log.Printf("failed get track list: %s", err)
		return err
	}

	log.Printf("fetched %d tracks from SoundCloud", len(likesList))

	log.Println("loading known tracks from database...")

	dbIDs, err := database.GetAllID()
	if err != nil {
		log.Printf("failed get track list GetAllID: %s", err)
		return err
	}

	log.Printf("found %d tracks already in database", len(dbIDs))

	existing := make(map[int64]struct{})
	for _, id := range dbIDs {
		existing[id] = struct{}{}
	}

	freshIDs := make(map[int64]struct{})
	for _, track := range likesList {
		freshIDs[track.ID] = struct{}{}
	}

	log.Println("checking for new tracks...")
	newCount := 0

	for _, track := range likesList {
		_, ok := existing[track.ID]
		if !ok {
			newCount++
			artist := track.Artist
			if artist == "" {
				artist = track.Uploader
			}
			log.Printf("[new %d] %s - %s", newCount, track.Title, artist)
			// все равно "- %s" пустой, хотя проверка выше есть
		}
	}

	log.Printf("total new tracks: %d", newCount)

	log.Println("checking for deleted tracks...")
	deletedCount := 0

	for id := range freshIDs {
		_, ok := freshIDs[id]
		if !ok {
			deletedCount++
			log.Printf("[deleted %d] track id: %d", deletedCount, id)
		}
	}

	log.Printf("total deleted tracks: %d", deletedCount)

	log.Println("=== sync-cloud: sync finished ===")

	return nil

}
