package categories

import "strings"

func Classify(ext string) string {
	ext = strings.ToLower(ext)
	photos := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".heic"}
	videos := []string{".mp4", ".mov", ".avi", ".mkv"}
	docs := []string{".pdf", ".txt", ".doc", ".docx", ".md"}
	compressed := []string{".zip", ".rar", ".7z", ".tar", ".gz", ".bz2", ".xz", ".tar.gz", ".tar.bz2", ".tar.xz", ".z", ".lzma", ".lz4", ".zst", ".arj", ".cab", ".deb", ".rpm", ".dmg", ".iso", ".img"}

	if contains(photos, ext) {
		return "photo"
	}

	if contains(videos, ext) {
		return "video"
	}

	if contains(docs, ext) {
		return "doc"
	}

	if contains(compressed, ext) {
		return "compressed"
	}

	return "other"
}

func contains(list []string, ext string) bool {
	for _, val := range list {
		if val == ext {
			return true
		}
	}
	return false
}
