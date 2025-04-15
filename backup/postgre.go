package backup

import (
	"compress/gzip"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

// PostgresBackup periodically backup postgres database in the home directory
func PostgresBackup(dirName string) {
	backupTime := getBackupTime()

	// new ticker to run every 6 hours
	ticker := time.NewTicker(backupTime)
	defer ticker.Stop()

	// Verify if backup directory is set
	if dirName == "" {
		fmt.Fprintln(os.Stderr, "Backup directory not set, We'll use home directory")
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to determine home directory, We'll use current directory")
			dirName = "."
		} else {
			dirName = homeDir
		}
	}

	// verify if backup directory exists
	if f, err := os.Stat(dirName); os.IsNotExist(err) || !f.IsDir() {
		fmt.Fprintf(os.Stderr, "Please create the backup directory '%s' first\n", dirName)
		os.Exit(1)
	}

	// Perform backup
	fmt.Fprintf(os.Stderr, "Performing backup every %s in %s\n", backupTime, dirName)
	for range ticker.C {
		err := performBackup(dirName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}

func performBackup(dirName string) error {
	// Retrieve database credentials from environment variables
	password, dbName, userName, host := getDBCredentials()
	fileName := createFileName()
	filePath := filepath.Join(dirName, fileName)

	// Perform backup using a secure approach
	if err := postgresDump(filePath, password, dbName, userName, host); err != nil {
		return fmt.Errorf("pg_dump error: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Backup completed successfully: %s\n", filePath)
	return nil
}

func postgresDump(fileName, password, dbName, userName, host string) error {
	// Set up the pg_dump command
	cmd := exec.Command("pg_dump", "-U", userName, "-d", dbName, "-h", host)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", password))

	// Create the output .gz file
	outFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create gzip file: %w", err)
	}
	defer outFile.Close()

	// Create a gzip writer
	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	// Pipe pg_dump stdout to gzip
	cmd.Stdout = gzipWriter

	// Run pg_dump
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pg_dump error: %w", err)
	}

	return nil
}

// getBackupTime get the backup time from the environment variable
func getBackupTime() time.Duration {
	defaultInterval := 6

	if os.Getenv("BACKUP_TIME") == "" {
		fmt.Fprintln(os.Stderr, "BACKUP_TIME env variable not set, using default value")
		return time.Duration(defaultInterval) * time.Hour
	}

	backupTime, err := strconv.Atoi(os.Getenv("BACKUP_TIME"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing BACKUP_TIME env variable, using default value")
		return time.Duration(defaultInterval) * time.Hour
	}

	if backupTime <= 0 || backupTime > 24 {
		fmt.Fprintln(os.Stderr, "Invalid BACKUP_TIME env variable, using default value")
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

func createFileName() string {
	var buf [32]byte
	t := time.Now()
	formatted := t.AppendFormat(buf[:0], "2006-01-02T15-04-05")
	return "backup_" + string(formatted) + ".sql.gz"
}
