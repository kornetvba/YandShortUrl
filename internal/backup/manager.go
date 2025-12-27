package backup

import "github.com/kornetvba/YandShortUrl/internal/config/config"

type FileExporter interface {
	SaveFile(filePath *config.FilePathType) error
	DownloadRecords(filePath *config.FilePathType) error
}

type BackupManager struct {
	Manager FileExporter
}

func NewBackupManager(manager FileExporter) *BackupManager {
	return &BackupManager{Manager: manager}
}

func (bm *BackupManager) SaveDataFile(filePath *config.FilePathType) error {
	return bm.Manager.SaveFile(filePath)
}

func (bm *BackupManager) DownloadRecordsToFile(filePath *config.FilePathType) error {
	return bm.Manager.DownloadRecords(filePath)
}
