package sync

import (
	"log"
	"os"
	"sync-cloud/internal/scrapper"
	"sync-cloud/internal/storage"
	"time"

	"gorm.io/gorm"
)

func Run(db *gorm.DB, profileLink, musicDir, proxyURL string) error {

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

	// existing - то, что СУЩЕСТВУЕТ у нас в бд
	existing := make(map[int64]struct{})
	for _, id := range dbIDs {
		existing[id] = struct{}{}
	}

	// freshIDs - актуальный список лайков на ск ПРЯМО СЕЙЧАС
	freshIDs := make(map[int64]struct{})
	for _, track := range likesList {
		freshIDs[track.ID] = struct{}{}
	}

	log.Println("checking for new tracks...")
	newCount := 0

	// идем по списку ЛАЙКНУТЫХ и сверяем с тем что у нас ЕСТЬ ПРЯМО СЕЙЧАС в бд
	// если трека в бд нет - скачиваем и добавляем
	for _, track := range likesList {
		_, ok := existing[track.ID]
		if !ok {
			newCount++

			time.Sleep(2 * time.Second)

			data, err := scrapper.GetMeta(track.URL)
			if err != nil {
				log.Printf("[error %d] %s get error: %s", newCount, track.Title, err)
				continue
			}

			outputPath, err := scrapper.DownloadAndTag(data, musicDir, proxyURL)
			if err != nil {
				log.Printf("[error %d] %s failed to download: %s", newCount, data.Title, err)
				continue
			}

			newTrack := &storage.LikedTrack{
				ID:       data.ID,
				Title:    data.Title,
				Artist:   data.Artist,
				FilePath: outputPath,
				LikedAt:  time.Now(),
			}

			if err := database.Create(newTrack); err != nil {
				log.Printf("[error %d] %s failed to save db: %s", newCount, data.Title, err)
				continue
			}

			log.Printf("[new %d] %s - %s", newCount, data.Title, data.Artist)
		}
	}

	log.Printf("total new tracks: %d", newCount)

	log.Println("checking for deleted tracks...")
	deletedCount := 0

	// идем по списку СУЩЕСТВУЮЗИХ В БД и сверяем с тем что ЛАЙКНУТО И АКТУАЛЬНО
	// если в лайкнутых трека нет - удаляем из бд
	for id := range existing {
		_, ok := freshIDs[id]
		if !ok {
			deletedCount++

			track, err := database.GetByID(id)
			if err != nil {
				log.Printf("[error %d] %d get by id error: %s", deletedCount, id, err)
				continue
			}

			os.Remove(track.FilePath)

			err = database.Delete(id)
			if err != nil {
				log.Printf("[error %d] %d get error: %s", deletedCount, id, err)
				continue
			}

			log.Printf("[deleted %d] track id: %d", deletedCount, id)
		}
	}

	log.Printf("total deleted tracks: %d", deletedCount)

	log.Println("=== sync-cloud: sync finished ===")

	return nil

}
