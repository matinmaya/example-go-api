package storagecache

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Cleaner struct {
	Dir    string
	MaxAge time.Duration
}

func NewCleaner(dir string, maxAge time.Duration) *Cleaner {
	return &Cleaner{
		Dir:    dir,
		MaxAge: maxAge,
	}
}

func (c *Cleaner) Start(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			c.ClearExpired()
		}
	}()
}

func (c *Cleaner) ClearExpired() {
	filepath.Walk(c.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("[storage-cache] %s", err.Error())
			return nil
		}
		if !info.IsDir() && time.Since(info.ModTime()) > c.MaxAge {
			_ = os.Remove(path)
			fmt.Println("[storage-cache] Removed expired:", path)
		}
		return nil
	})

	var dirs []string
	filepath.Walk(c.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("[storage-cache] %s", err.Error())
			return nil
		}
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	for i := len(dirs) - 1; i >= 0; i-- {
		if dirs[i] == c.Dir {
			continue
		}
		entries, err := os.ReadDir(dirs[i])
		if err == nil && len(entries) == 0 {
			_ = os.Remove(dirs[i])
			fmt.Println("[storage-cache] Removed empty folder:", dirs[i])
		}
	}
}
