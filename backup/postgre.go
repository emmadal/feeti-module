package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

// PostgresBackup periodically backup postgres database in the home directory
func PostgresBackup(dirName string) {
	ticker := time.NewTicker(10 * time.Second) // new ticker to run every 12 hours
	defer ticker.Stop()

	// loop run every 10 seconds
	for range ticker.C {
		err := performBackup(dirName)
		if err != nil {
			logrus.WithFields(logrus.Fields{"error": err}).Error(err.Error())
		}
	}
}

func performBackup(dirName string) error {
	// create backup file name and open file
	fileName, file, err := createFileName()
	if err != nil {
		return fmt.Errorf("backup file creation error: %w", err)
	}
	defer file.Close()

	// perform backup
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	userName := os.Getenv("DB_USER")
	host := os.Getenv("DB_HOST")

	// perform dump
	err = postgresDump(fileName, password, dbName, userName, host)
	if err != nil {
		return fmt.Errorf("dump error: %w", err)
	}

	// create backup directory
	directory, err := createBackupDir(dirName)
	if err != nil {
		return fmt.Errorf("backup directory creation error: %w", err)
	}

	// move file to backup directory
	err = os.Rename(fileName, filepath.Join(directory, fileName))
	if err != nil {
		return fmt.Errorf("backup rename error: %w", err)
	}

	logrus.Infof("Backup completed successfully: %s", fileName)
	return nil
}

func createFileName() (string, *os.File, error) {
	fileName := fmt.Sprintf("backup_%s.gz", time.Now().Format("2006-01-02T15:04:05"))
	file, err := os.Create(fileName)
	if err != nil {
		return "", nil, fmt.Errorf("backup file error: %w", err)
	}
	return fileName, file, nil
}

func postgresDump(fileName, password, dbName, userName, host string) error {
	query := fmt.Sprintf("PGPASSWORD='%s' pg_dump -U %s -d %s -h %s | gzip > %s",
		password, userName, dbName, host, fileName)
	cmd := exec.Command("sh", "-c", query)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("pg_dump error: %w", err)
	}
	return nil
}

func createBackupDir(dirName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home directory error: %w", err)
	}

	directory := filepath.Join(homeDir, dirName)

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err = os.MkdirAll(directory, 0755)
		if err != nil {
			return "", fmt.Errorf("backup directory creation error: %w", err)
		}
	}
	return directory, nil
}
