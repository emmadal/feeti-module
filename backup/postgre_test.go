package backup

import (
	"os"
	"testing"
)

func BenchmarkGetBackupTime(b *testing.B) {
	// Save the original env var to restore later
	originalBackupTime := os.Getenv("BACKUP_TIME")
	defer func() {
		_ = os.Setenv("BACKUP_TIME", originalBackupTime)
	}()

	// Set environment variable for testing
	_ = os.Setenv("BACKUP_TIME", "12")

	b.ReportAllocs()

	for b.Loop() {
		_ = getBackupTime()
	}
}

func BenchmarkGetDBCredentials(b *testing.B) {
	// Save original env vars to restore later
	originalPassword := os.Getenv("DB_PASSWORD")
	originalDBName := os.Getenv("DB_NAME")
	originalDBUser := os.Getenv("DB_USER")
	originalDBHost := os.Getenv("DB_HOST")
	defer func() {
		_ = os.Setenv("DB_PASSWORD", originalPassword)
		_ = os.Setenv("DB_NAME", originalDBName)
		_ = os.Setenv("DB_USER", originalDBUser)
		_ = os.Setenv("DB_HOST", originalDBHost)
	}()

	// Set environment variables for testing
	_ = os.Setenv("DB_PASSWORD", "test_password")
	_ = os.Setenv("DB_NAME", "test_db")
	_ = os.Setenv("DB_USER", "test_user")
	_ = os.Setenv("DB_HOST", "localhost")

	b.ReportAllocs()

	for b.Loop() {
		_, _, _, _ = getDBCredentials()
	}
}

// BenchmarkFileName benchmarks the generation of the backup filename
func BenchmarkFileName(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = createFileName()
	}
}

// BenchmarkPerformBackup benchmarks the core parts of the backup process
// without actually executing pg_dump
func BenchmarkPerformBackup(b *testing.B) {
	// Create a test directory
	testDir := b.TempDir()

	// Set required environment variables
	originalEnv := setupTestEnv()
	defer restoreEnv(originalEnv)

	b.ReportAllocs()

	for b.Loop() {
		_ = performBackup(testDir)
	}
}

// Helper function to set up test environment and save original values
func setupTestEnv() map[string]string {
	originalEnv := make(map[string]string)

	envVars := []string{"DB_PASSWORD", "DB_NAME", "DB_USER", "DB_HOST"}
	for _, v := range envVars {
		originalEnv[v] = os.Getenv(v)
	}

	_ = os.Setenv("DB_PASSWORD", "test_password")
	_ = os.Setenv("DB_NAME", "test_db")
	_ = os.Setenv("DB_USER", "test_user")
	_ = os.Setenv("DB_HOST", "localhost")

	return originalEnv
}

// Helper function to restore original environment
func restoreEnv(originalEnv map[string]string) {
	for k, v := range originalEnv {
		_ = os.Setenv(k, v)
	}
}
