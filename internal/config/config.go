package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	ProfileURL string
	DBPath     string
	MusicDir   string
	ProxyURL   string
}

func Load() (*Config, error) {
	godotenv.Load()

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		defaultPath, err := defaultDBPath()
		if err != nil {
			return nil, err
		}
		dbPath = defaultPath
	}

	profileURL := os.Getenv("PROFILE_URL")
	if profileURL == "" {
		return nil, fmt.Errorf("PROFILE_URL is not set")
	}

	musicDir := os.Getenv("MUSIC_DIR")
	if musicDir == "" {
		defaultMusicDir, err := defaultMusicDir()
		if err != nil {
			return nil, err
		}
		musicDir = defaultMusicDir
	}

	proxyURL := os.Getenv("PROXY_URL")

	return &Config{
		ProfileURL: profileURL,
		DBPath:     dbPath,
		MusicDir:   musicDir,
		ProxyURL:   proxyURL,
	}, nil

}

func defaultDBPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	exeDir := filepath.Dir(exePath)
	dbDir := filepath.Join(exeDir, "data", "database")

	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create db directory: %w", err)
	}

	return filepath.Join(dbDir, "likes.db"), nil
}

func defaultMusicDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	exeDir := filepath.Dir(exePath)
	dbDir := filepath.Join(exeDir, "data", "music")

	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create db directory: %w", err)
	}

	return dbDir, nil
}
