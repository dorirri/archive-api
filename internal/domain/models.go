package domain

type FileInfo struct {
	FilePath string  `json:"file_path"`
	Size     float64 `json:"size"`
	Mimetype string  `json:"mimetype"`
}

type ArchiveInfo struct {
	FileName    string     `json:"filename"`
	ArchiveSize float64    `json:"archive_size"`
	TotalSize   float64    `json:"total_size"`
	TotalFiles  float64    `json:"total_files"`
	Files       []FileInfo `json:"files"`
}

type ArchiveRepository interface {
	AnalyzeArchive(content []byte, filename string) (*ArchiveInfo, error)
}

type AllowedMimeTypes struct {
	ArchiveTypes map[string]bool
	EmailTypes   map[string]bool
}

var ALlowedTypes = AllowedMimeTypes{
	ArchiveTypes: map[string]bool{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/xml": true,
		"image/jpeg":      true,
		"image/png":       true,
	},
	EmailTypes: map[string]bool{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/pdf": true,
	},
}
