package scanner

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/SanujaRubasinghe/ssdindexer/internal/categories"
)

func Scan(root string, updateChan chan<- FileStats) (FileStats, error) {
	var stats FileStats
	stats.OtherExtensions = make(map[string]int64)
	folderSizes := make(map[string]int64)
	var folderMutex sync.Mutex

	// First pass: count total files
	totalFiles := countFiles(root)
	if totalFiles == 0 {
		return stats, nil
	}

	workerCount := runtime.NumCPU()
	files := make(chan string, 10000)
	results := make(chan FileStats, workerCount)
	progressChan := make(chan int64, workerCount)

	// Start workers
	for i := 0; i < workerCount; i++ {
		go func() {
			var local FileStats
			local.OtherExtensions = make(map[string]int64)
			var processed int64
			for path := range files {
				info, err := os.Stat(path)
				if err != nil || info.IsDir() {
					continue
				}

				size := info.Size()
				ext := filepath.Ext(path)
				category := categories.Classify(ext)
				
				switch category {
				case "photo":
					local.Photos += size
				case "video":
					local.Videos += size
				case "doc":
					local.Docs += size
				case "compressed":
					local.Compressed += size
				default:
					local.Others += size
					// Track extension details only for Other category
					local.OtherExtensions[ext] += size
				}
				local.Total += size

				// Track folder sizes
				dir := filepath.Dir(path)
				folderMutex.Lock()
				folderSizes[dir] += size
				folderMutex.Unlock()
				processed++

				// Send progress update more frequently for better UX
				if processed%5 == 0 {
					progressPercent := float64(processed) / float64(totalFiles) * 1000000.0
					updateChan <- FileStats{
						Photos:     local.Photos,
						Videos:     local.Videos,
						Docs:       local.Docs,
						Compressed: local.Compressed,
						Others:     local.Others,
						Total:      int64(progressPercent),
						OtherExtensions: local.OtherExtensions,
					}
				}
			}
			results <- local
			progressChan <- processed
		}()
	}

	// File discovery goroutine
	go func() {
		filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err == nil && !d.IsDir() {
				files <- path
			}
			return nil
		})
		close(files)
	}()

	// Collect results
	var totalProcessed int64
	for i := 0; i < workerCount; i++ {
		r := <-results
		stats.Photos += r.Photos
		stats.Videos += r.Videos
		stats.Docs += r.Docs
		stats.Compressed += r.Compressed
		stats.Others += r.Others
		stats.Total += r.Total
		
		// Merge Other extension maps
		for ext, size := range r.OtherExtensions {
			stats.OtherExtensions[ext] += size
		}
		
		totalProcessed += <-progressChan
	}

	// Convert folder sizes to slice for sorting
	var folderSizesSlice []FolderSize
	for path, size := range folderSizes {
		folderSizesSlice = append(folderSizesSlice, FolderSize{
			Path: path,
			Size: size,
		})
	}

	// Send final update with 100% progress and folder sizes
	stats.FolderSizes = folderSizesSlice
	updateChan <- FileStats{
		Photos:         stats.Photos,
		Videos:         stats.Videos,
		Docs:           stats.Docs,
		Compressed:     stats.Compressed,
		Others:         stats.Others,
		Total:          1000000, // 100% complete
		OtherExtensions: stats.OtherExtensions,
		FolderSizes:    folderSizesSlice,
	}

	return stats, nil
}

func countFiles(root string) int64 {
	var count int64
	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			count++
		}
		return nil
	})
	return count
}
