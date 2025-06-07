// فایل: cmd/backup_mongo.go
package cmd

import (
	// ... سایر import ها مثل قبل ...
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// نام متغیر به backupMongoCmd تغییر کرد
var backupMongoCmd = &cobra.Command{
	Use:   "mongo",
	Short: "Backup a MongoDB database",
	Long: `Creates a compressed archive of a remote MongoDB database using mongodump.
It reads the required configurations (remote URI, tool path, backup path)
from the application's config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// ... محتوای تابع Run دقیقاً مثل قبل باقی می‌ماند ...
		fmt.Println("Starting MongoDB backup...")

		remoteURI := viper.GetString("mongodb.remote_uri")
		toolsPath := viper.GetString("paths.mongo_tools")
		backupDir := viper.GetString("paths.backup")

		if remoteURI == "" || toolsPath == "" || backupDir == "" {
			log.Fatal("Configuration error: 'mongodb.remote_uri', 'paths.mongo_tools', and 'paths.backup' must be set. Please run 'dataweaver-cli configure' first.")
		}

		fmt.Printf("Using remote URI: %s\n", remoteURI)
		fmt.Printf("Path to MongoDB tools: %s\n", toolsPath)
		fmt.Printf("Backup destination directory: %s\n", backupDir)

		mongoDumpExecutable := "mongodump"
		if runtime.GOOS == "windows" {
			mongoDumpExecutable = "mongodump.exe"
		}
		mongoDumpPath := filepath.Join(toolsPath, mongoDumpExecutable)

		if _, err := os.Stat(mongoDumpPath); os.IsNotExist(err) {
			log.Fatalf("mongodump not found at the specified path: %s. Please verify your 'paths.mongo_tools' configuration.", mongoDumpPath)
		}

		timestamp := time.Now().Format("2006-01-02_15-04-05")
		backupFileName := fmt.Sprintf("backup-%s.gz", timestamp)
		if err := os.MkdirAll(backupDir, 0755); err != nil {
			log.Fatalf("Failed to create backup directory '%s': %v", backupDir, err)
		}
		backupFilePath := filepath.Join(backupDir, backupFileName)
		fmt.Printf("Backup file will be saved to: %s\n", backupFilePath)

		dumpCmd := exec.Command(mongoDumpPath,
			fmt.Sprintf("--uri=%s", remoteURI),
			fmt.Sprintf("--archive=%s", backupFilePath),
			"--gzip",
		)

		stdout, err := dumpCmd.StdoutPipe()
		if err != nil {
			log.Fatalf("Failed to get stdout pipe: %v", err)
		}
		dumpCmd.Stderr = dumpCmd.Stdout

		scanner := bufio.NewScanner(stdout)
		go func() {
			for scanner.Scan() {
				fmt.Printf("  [mongodump]: %s\n", scanner.Text())
			}
		}()

		fmt.Println("Executing mongodump command...")
		if err := dumpCmd.Start(); err != nil {
			log.Fatalf("Failed to start mongodump command: %v", err)
		}

		err = dumpCmd.Wait()
		if err != nil {
			log.Fatalf("mongodump command failed with error: %v", err)
		}

		fmt.Println("------------------------")
		fmt.Println("MongoDB backup completed successfully!")
	},
}

func init() {
	// این دستور، خودش را به والدش (backupCmd) اضافه می‌کند
	backupCmd.AddCommand(backupMongoCmd)
}
