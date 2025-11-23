package scanner

import "sort"

type ExtensionStats struct {
	Ext  string
	Size int64
}

type FolderSize struct {
	Path string
	Size int64
}

type FileStats struct {
	Photos     int64
	Videos     int64
	Docs       int64
	Compressed int64
	Others     int64
	Total      int64
	OtherExtensions map[string]int64
	FolderSizes []FolderSize
}

func (fs *FileStats) GetTopOtherExtensions(limit int) []ExtensionStats {
	var exts []ExtensionStats
	for ext, size := range fs.OtherExtensions {
		exts = append(exts, ExtensionStats{Ext: ext, Size: size})
	}
	
	sort.Slice(exts, func(i, j int) bool {
		return exts[i].Size > exts[j].Size
	})
	
	if len(exts) > limit {
		exts = exts[:limit]
	}
	
	return exts
}

type ProgressMsg struct {
	Current int64
	Total   int64
	Stats   FileStats
}

type DoneMsg struct{}

func (fs *FileStats) GetTopFolders(limit int) []FolderSize {
	// Make a copy of the folder sizes to avoid modifying the original
	folders := make([]FolderSize, len(fs.FolderSizes))
	copy(folders, fs.FolderSizes)
	
	sort.Slice(folders, func(i, j int) bool {
		return folders[i].Size > folders[j].Size
	})
	
	if len(folders) > limit {
		return folders[:limit]
	}
	return folders
}
