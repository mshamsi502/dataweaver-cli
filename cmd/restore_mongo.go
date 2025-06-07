// فایل: cmd/restore_mongo.go
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// restoreMongoCmd represents the mongo subcommand of restore
var restoreMongoCmd = &cobra.Command{
	Use:   "mongo",
	Short: "Restore a MongoDB database from a backup file",
	Long: `Restores a MongoDB database from a previously created archive file.
You will be prompted to select a backup file from the configured backup directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting MongoDB restore...")

		// 1. خواندن تنظیمات
		localURI := viper.GetString("mongodb.local_uri")
		toolsPath := viper.GetString("paths.mongo_tools")
		backupDir := viper.GetString("paths.backup")

		// 2. بررسی تنظیمات ضروری
		if localURI == "" || toolsPath == "" || backupDir == "" {
			log.Fatal("Configuration error: 'mongodb.local_uri', 'paths.mongo_tools', and 'paths.backup' must be set.")
		}

		// 3. پیدا کردن فایل‌های بکاپ موجود
		backupFiles, err := findBackupFiles(backupDir)
		if err != nil {
			log.Fatalf("Error reading backup directory '%s': %v", backupDir, err)
		}
		if len(backupFiles) == 0 {
			log.Fatalf("No backup files found in '%s'.", backupDir)
		}

		// 4. گرفتن ورودی از کاربر برای انتخاب فایل بکاپ
		var selectedFile string
		prompt := &survey.Select{
			Message: "Choose a backup to restore:",
			Options: backupFiles,
		}
		survey.AskOne(prompt, &selectedFile, survey.WithValidator(survey.Required))

		backupFilePath := filepath.Join(backupDir, selectedFile)
		fmt.Printf("Selected backup file: %s\n", backupFilePath)

		// 5. ساختن و اجرای دستور mongorestore
		mongoRestoreExecutable := "mongorestore"
		if runtime.GOOS == "windows" {
			mongoRestoreExecutable += ".exe"
		}
		mongoRestorePath := filepath.Join(toolsPath, mongoRestoreExecutable)

		if _, err := os.Stat(mongoRestorePath); os.IsNotExist(err) {
			log.Fatalf("mongorestore not found at the specified path: %s. Please verify your 'paths.mongo_tools' configuration.", mongoRestorePath)
		}

		// The --drop flag will drop collections from the target database before restoring.
		restoreCmd := exec.Command(mongoRestorePath,
			fmt.Sprintf("--uri=%s", localURI),
			fmt.Sprintf("--archive=%s", backupFilePath),
			"--gzip",
			"--drop",
		)

		fmt.Println("Executing mongorestore command. This might take a while...")

		// اجرای دستور و نمایش خروجی به صورت زنده
		restoreCmd.Stdout = os.Stdout
		restoreCmd.Stderr = os.Stderr

		if err := restoreCmd.Run(); err != nil {
			log.Fatalf("mongorestore command failed: %v", err)
		}

		fmt.Println("------------------------")
		fmt.Println("MongoDB restore completed successfully!")
	},
}

// findBackupFiles finds all .gz backup files in the specified directory
func findBackupFiles(backupDir string) ([]string, error) {
	var files []string
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".gz") {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

func init() {
	// اضافه کردن این زیردستور به دستور والد 'restore'
	restoreCmd.AddCommand(restoreMongoCmd)
}
