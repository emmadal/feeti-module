package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// PostgresBackup periodically backup postgres database in the home directory
func PostgresBackup(dirName string) {
	backupTime := getBackupTime()

	// new ticker to run every 6 hours
	ticker := time.NewTicker(backupTime)
	defer ticker.Stop()

	// Verify if backup directory is set
	if dirName == "" {
		logrus.Warn("Backup directory not set, We'll use home directory")
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logrus.Warn("Failed to determine home directory, We'll use current directory")
			dirName = "."
		} else {
			dirName = homeDir
		}
	}

	// verify if backup directory exists
	if f, err := os.Stat(dirName); os.IsNotExist(err) || !f.IsDir() {
		logrus.Errorf("Please create the backup directory '%s' first", dirName)
		os.Exit(1)
	}

	// Perform backup
	logrus.Infof("Performing backup every %s in %s", backupTime, dirName)
	for range ticker.C {
		err := performBackup(dirName)
		if err != nil {
			logrus.WithFields(logrus.Fields{"error": err}).Error(err.Error())
		}
	}
}

func performBackup(dirName string) error {
	start := time.Now()

	// Retrieve database credentials from environment variables
	password, dbName, userName, host := getDBCredentials()
	fileName := fmt.Sprintf("backup_%s.sql.gz", time.Now().Format("2006-01-02T15-04-05"))
	filePath := filepath.Join(dirName, fileName)

	// Perform backup using a secure approach
	if err := postgresDump(filePath, password, dbName, userName, host); err != nil {
		return fmt.Errorf("pg_dump error: %w", err)
	}

	logrus.Infof("Backup completed successfully: %s", filePath)
	logrus.Infof("Backup took %s", time.Since(start))
	return nil
}

func postgresDump(fileName, password, dbName, userName, host string) error {
	// perform backup
	query := fmt.Sprintf("pg_dump -U %s -d %s -h %s | gzip > %s", userName, dbName, host, fileName)
	cmd := exec.Command("sh", "-c", query)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", password)) // set password in the environment

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()

	cmd.Stdout = file

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("pg_dump error: %w", err)
	}
	return nil
}

// getBackupTime get the backup time from the environment variable
func getBackupTime() time.Duration {
	defaultInterval := 6

	if os.Getenv("BACKUP_TIME") == "" {
		logrus.Warningln("BACKUP_TIME env variable not set, using default value")
		return time.Duration(defaultInterval) * time.Hour
	}

	backupTime, err := strconv.Atoi(os.Getenv("BACKUP_TIME"))
	if err != nil {
		logrus.Warningln("Error parsing BACKUP_TIME env variable, using default value")
		return time.Duration(defaultInterval) * time.Hour
	}

	if backupTime <= 0 || backupTime > 24 {
		logrus.Warningln("Invalid BACKUP_TIME env variable, using default value")
		return time.Duration(defaultInterval) * time.Hour
	}
	return time.Duration(backupTime) * time.Hour
}

// getDBCredentials get the database credentials from the environment variables
func getDBCredentials() (string, string, string, string) {
	return os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_HOST")
}
