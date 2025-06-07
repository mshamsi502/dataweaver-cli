package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var downloadToolsCmd = &cobra.Command{
	Use:   "download-tools",
	Short: "Downloads and extracts MongoDB Database Tools into a local directory.",
	Long: `This command automates the setup of MongoDB Database Tools for portable use.
It downloads the correct ZIP archive for the OS, extracts it into a local 'tools' directory,
and updates the configuration with the path to the executables.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting MongoDB Database Tools setup...")

		// --- 1. تعیین URL و مسیرهای لازم ---
		// این بار از لینک فایل ZIP استفاده می‌کنیم
		// directDownloadURL := "https://fastdl.mongodb.org/tools/db/mongodb-database-tools-windows-x86_64-100.12.1.zip"
		directDownloadURL := "https://fastdl.mongodb.org/tools/db/mongodb-database-tools-windows-x86_64-100.12.2.zip"

		// TODO: در آینده این بخش می‌تواند بر اساس سیستم‌عامل، URL مناسب را انتخاب کند.
		if runtime.GOOS != "windows" {
			log.Fatal("This version currently only supports downloading the Windows ZIP archive.")
		}

		fileName := filepath.Base(directDownloadURL)
		downloadDir := "./downloads/"
		extractDir := "./tools/" // پوشه‌ای که ابزارها در آن استخراج می‌شوند
		downloadFilePath := filepath.Join(downloadDir, fileName)

		// --- 2. دانلود فایل ZIP (اگر وجود ندارد) ---
		fmt.Printf("Checking for installer at: %s\n", downloadFilePath)
		if err := os.MkdirAll(downloadDir, 0755); err != nil {
			log.Fatalf("Error creating download directory %s: %v", downloadDir, err)
		}

		if _, err := os.Stat(downloadFilePath); os.IsNotExist(err) {
			fmt.Println("Installer not found. Downloading...")
			downloadFile(directDownloadURL, downloadFilePath)
		} else {
			fmt.Println("Installer ZIP file already exists. Skipping download.")
		}

		// --- 3. استخراج فایل ZIP ---
		fmt.Printf("Extracting '%s' to '%s'...\n", downloadFilePath, extractDir)
		binPath, err := unzipAndFindBin(downloadFilePath, extractDir)
		if err != nil {
			log.Fatalf("Failed to extract and find tools: %v", err)
		}
		fmt.Println("Extraction complete.")

		// --- 4. آپدیت کردن کانفیگ با مسیر نهایی ---
		updateMongoToolsPathInConfig(binPath)
	},
}

// تابع کمکی برای آپدیت کانفیگ (کمی تمیزتر شده)
func updateMongoToolsPathInConfig(toolPath string) {
	absToolPath, err := filepath.Abs(toolPath)
	if err != nil {
		log.Printf("Warning: could not get absolute path for '%s'. Storing relative path.", toolPath)
		absToolPath = toolPath // در صورت خطا، از همان مسیر نسبی استفاده می‌کنیم
	}

	fmt.Printf("Updating configuration: 'paths.mongo_tools' -> '%s'\n", absToolPath)
	viper.Set("paths.mongo_tools", absToolPath)

	// پیدا کردن مسیر فایل کانفیگ برای ذخیره
	home, _ := os.UserHomeDir()
	configFile := filepath.Join(home, ".dataweaver-cli", "config.yaml")

	// اگر مسیر خانگی در دسترس نیست یا پوشه وجود ندارد، در پوشه فعلی ذخیره کن
	if _, err := os.Stat(filepath.Dir(configFile)); os.IsNotExist(err) {
		configFile = "config.yaml"
	}

	if err := viper.WriteConfigAs(configFile); err != nil {
		log.Printf("Error writing configuration to '%s': %v\n", configFile, err)
	} else {
		fmt.Println("-------------------------------------------------")
		fmt.Printf("SUCCESS: Configuration updated in '%s'\n", configFile)
		fmt.Println("-------------------------------------------------")
	}
}

// تابع کمکی برای دانلود فایل (بدون تغییر)
func downloadFile(url, destPath string) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Download error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Download failed with status: %s", resp.Status)
	}
	outFile, err := os.Create(destPath)
	if err != nil {
		log.Fatalf("Cannot create file %s: %v", destPath, err)
	}
	defer outFile.Close()
	size, err := io.Copy(outFile, resp.Body)
	if err != nil {
		os.Remove(destPath)
		log.Fatalf("Error saving file: %v", err)
	}
	fmt.Printf("Successfully downloaded '%s' (%d bytes).\n", filepath.Base(destPath), size)
}

// تابع جدید برای استخراج فایل ZIP و پیدا کردن مسیر پوشه bin
func unzipAndFindBin(source, destination string) (string, error) {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	var binPath string
	foundBin := false

	for _, file := range reader.File {
		// مسیر کامل فایل در مقصد
		fpath := filepath.Join(destination, file.Name)

		// جلوگیری از حملات Zip Slip
		if !strings.HasPrefix(fpath, filepath.Clean(destination)+string(os.PathSeparator)) {
			return "", fmt.Errorf("illegal file path: %s", fpath)
		}

		// اگر ورودی یک پوشه است، آن را ایجاد کن
		if file.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// اگر فایل است، ابتدا پوشه والد آن را ایجاد کن
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return "", err
		}

		// فایل را باز کن
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return "", err
		}

		rc, err := file.Open()
		if err != nil {
			return "", err
		}

		_, err = io.Copy(outFile, rc)

		// منابع را ببند
		outFile.Close()
		rc.Close()

		if err != nil {
			return "", err
		}

		// بررسی برای پیدا کردن پوشه bin
		if !foundBin && strings.HasSuffix(strings.ToLower(fpath), "mongodump.exe") {
			binPath = filepath.Dir(fpath)
			foundBin = true
		}
	}

	if !foundBin {
		return "", fmt.Errorf("could not find 'bin' directory containing 'mongodump.exe' in the zip file")
	}

	return binPath, nil
}

func init() {
	rootCmd.AddCommand(downloadToolsCmd)
}
